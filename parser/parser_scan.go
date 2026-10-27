package parser

import (
	"strings"

	"github.com/czerwonk/bird_exporter/protocol"
)

// scanMatchers match lines by scanning bytes instead of running regular
// expressions. Each matcher mirrors the regular expression of the same name in
// parser_re.go and returns the same submatches for every input.
var scanMatchers = lineMatchers{
	protocol:    scanProtocol,
	description: scanDescription,
	message:     scanMessage,
	channel:     scanChannel,
	routes:      scanRoutes,
	routeChange: scanRouteChange,
	filter:      scanFilter,
}

// Fixed fragments of `show protocols all` reply lines.
const (
	protocolReplyCode = "1002-"
	descriptionKey    = "Description:"
	messageKey        = "Message:"
	channelKey        = "Channel "
	channelIPv        = "ipv"
	channelROA        = "roa"
	routesKey         = "Routes:"
	listSeparator     = ", "
)

// ParseProtocolsScanWithError parses BIRD output like ParseProtocolsWithError,
// but without regular expressions.
func ParseProtocolsScanWithError(data []byte, ipVersion string) ([]*protocol.Protocol, error) {
	return parseProtocols(data, ipVersion, &scanMatchers)
}

// isSpace matches the \s class of Go regular expressions.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\f' || b == '\r'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func skipSpace(s string, i int) int {
	for i < len(s) && isSpace(s[i]) {
		i++
	}
	return i
}

func skipNonSpace(s string, i int) int {
	for i < len(s) && !isSpace(s[i]) {
		i++
	}
	return i
}

func skipDigits(s string, i int) int {
	for i < len(s) && isDigit(s[i]) {
		i++
	}
	return i
}

// leadingSpace returns the index after a leading whitespace run, or -1 when
// the line does not start with whitespace.
func leadingSpace(s string) int {
	i := skipSpace(s, 0)
	if i == 0 {
		return -1
	}
	return i
}

// word returns the non-whitespace run at i and the index after the whitespace
// run that must follow it.
func word(s string, i int) (string, int, bool) {
	end := skipNonSpace(s, i)
	if end == i {
		return "", i, false
	}

	next := skipSpace(s, end)
	if next == end {
		return "", i, false
	}

	return s[i:end], next, true
}

func isProtoName(s string) bool {
	switch s {
	case "MRT", "BGP", "BFD", "OSPF", "RPKI", "RIP", "RAdv", "Pipe", "Perf", "Direct", "Babel", "Device", "Kernel", "Static":
		return true
	}
	return false
}

// scanProtocol mirrors protocolRegex.
func scanProtocol(line string) []string {
	if strings.HasPrefix(line, protocolReplyCode) {
		if m := scanProtocolAt(line, len(protocolReplyCode)); m != nil {
			return m
		}
	}

	return scanProtocolAt(line, 0)
}

func scanProtocolAt(line string, i int) []string {
	name, i, ok := word(line, i)
	if !ok {
		return nil
	}

	proto, i, ok := word(line, i)
	if !ok || !isProtoName(proto) {
		return nil
	}

	table, i, ok := word(line, i)
	if !ok {
		return nil
	}

	state, i, ok := word(line, i)
	if !ok || i == len(line) {
		return nil
	}

	end := scanSince(line, i)
	info := ""
	if end < len(line) {
		info = line[skipSpace(line, end):]
	}

	return []string{line, name, proto, table, state, line[i:end], info}
}

// scanSince returns the end of the since column: a date and time with
// optional fraction when followed by whitespace or the end of the line,
// otherwise the non-whitespace run.
func scanSince(line string, i int) int {
	if end := scanDateTime(line, i); end >= 0 {
		if end == len(line) || isSpace(line[end]) {
			return end
		}

		if line[end] == '.' {
			frac := skipDigits(line, end+1)
			if frac > end+1 && (frac == len(line) || isSpace(line[frac])) {
				return frac
			}
		}
	}

	return skipNonSpace(line, i)
}

// scanDateTime matches "dddd-dd-dd dd:dd:dd" at i and returns its end, or -1.
func scanDateTime(s string, i int) int {
	const layout = "dddd-dd-dd dd:dd:dd"
	if len(s)-i < len(layout) {
		return -1
	}

	for j := 0; j < len(layout); j++ {
		c := s[i+j]
		if layout[j] == 'd' {
			if !isDigit(c) {
				return -1
			}
		} else if c != layout[j] {
			return -1
		}
	}

	return i + len(layout)
}

// afterReplyCode strips a leading "1006-" reply code like the optional group
// in descriptionRegex, messageRegex, routeChangeRegex and filterRegex.
func afterReplyCode(line string) string {
	return strings.TrimPrefix(line, "1006-")
}

// scanDescription mirrors descriptionRegex.
func scanDescription(line string) []string {
	s := afterReplyCode(line)
	i := leadingSpace(s)
	if i < 0 || !strings.HasPrefix(s[i:], descriptionKey) {
		return nil
	}

	i += len(descriptionKey)
	value := skipSpace(s, i)
	if value == i {
		return nil
	}

	return []string{line, s[value:]}
}

// scanMessage mirrors messageRegex.
func scanMessage(line string) []string {
	s := afterReplyCode(line)
	i := leadingSpace(s)
	if i < 0 || !strings.HasPrefix(s[i:], messageKey) {
		return nil
	}

	end := len(line) - len(s) + i + len(messageKey)
	return []string{line[:end]}
}

// scanChannel mirrors channelRegex.
func scanChannel(line string) []string {
	i := leadingSpace(line)
	if i < 0 || !strings.HasPrefix(line[i:], channelKey) {
		return nil
	}

	i += len(channelKey)
	if !strings.HasPrefix(line[i:], channelIPv) && !strings.HasPrefix(line[i:], channelROA) {
		return nil
	}

	i += len(channelIPv) // same length as channelROA
	if i == len(line) || (line[i] != '4' && line[i] != '6') {
		return nil
	}

	end := i + 1
	if end < len(line) {
		if !isSpace(line[end]) {
			return nil
		}
		end++
	}

	return []string{line[:end], line[i : i+1]}
}

// scanRoutes mirrors routeRegex.
func scanRoutes(line string) []string {
	i := leadingSpace(line)
	if i < 0 || !strings.HasPrefix(line[i:], routesKey) {
		return nil
	}

	i += len(routesKey)
	start := skipSpace(line, i)
	if start == i {
		return nil
	}

	imported, i, ok := numberFollowedBy(line, start, " imported, ")
	if !ok {
		return nil
	}

	filtered := ""
	if n, next, ok := numberFollowedBy(line, i, " filtered, "); ok {
		filtered, i = n, next
	}

	exported, i, ok := numberFollowedBy(line, i, " exported")
	if !ok {
		return nil
	}

	preferred := ""
	if strings.HasPrefix(line[i:], listSeparator) {
		if n, next, ok := numberFollowedBy(line, i+len(listSeparator), " preferred"); ok {
			preferred, i = n, next
		}
	}

	return []string{line[:i], imported, filtered, exported, preferred}
}

// numberFollowedBy matches a digit run at i followed by suffix.
func numberFollowedBy(s string, i int, suffix string) (string, int, bool) {
	end := skipDigits(s, i)
	if end == i || !strings.HasPrefix(s[end:], suffix) {
		return "", i, false
	}

	return s[i:end], end + len(suffix), true
}

// Submatch layout of routeChangeRegex: full match, direction, kind, the
// required counters, then the two optional v3 counters (rx limit, limit).
const (
	routeChangeGroups       = 9
	routeChangeFirstCounter = 3
	routeChangeCounters     = 5
)

// scanRouteChange mirrors routeChangeRegex.
func scanRouteChange(line string) []string {
	s := afterReplyCode(line)
	i := leadingSpace(s)
	if i < 0 {
		return nil
	}

	var direction string
	switch {
	case strings.HasPrefix(s[i:], routeChangeImport+" "):
		direction = routeChangeImport
	case strings.HasPrefix(s[i:], routeChangeExport+" "):
		direction = routeChangeExport
	default:
		return nil
	}
	i += len(direction) + 1

	var kind string
	switch {
	case strings.HasPrefix(s[i:], routeChangeUpdates+":"):
		kind = routeChangeUpdates
	case strings.HasPrefix(s[i:], routeChangeWithdraws+":"):
		kind = routeChangeWithdraws
	default:
		return nil
	}
	i += len(kind) + 1

	m := make([]string, 1+routeChangeGroups)
	m[1], m[2] = direction, kind

	for n := routeChangeFirstCounter; n < routeChangeFirstCounter+routeChangeCounters; n++ {
		value, next, ok := spacedRouteChangeValue(s, i)
		if !ok {
			return nil
		}
		m[n], i = value, next
	}

	if rxLimit, next, ok := spacedRouteChangeValue(s, i); ok {
		if limit, end, ok := spacedRouteChangeValue(s, next); ok {
			m[routeChangeGroups-1], m[routeChangeGroups] = rxLimit, limit
			i = end
		}
	}

	m[0] = line[:len(line)-len(s)+i]
	return m
}

// spacedRouteChangeValue matches whitespace followed by a digit run or "---".
func spacedRouteChangeValue(s string, i int) (string, int, bool) {
	start := skipSpace(s, i)
	if start == i {
		return "", i, false
	}

	if end := skipDigits(s, start); end > start {
		return s[start:end], end, true
	}

	if strings.HasPrefix(s[start:], routeChangeUnset) {
		return routeChangeUnset, start + len(routeChangeUnset), true
	}

	return "", i, false
}

// scanFilter mirrors filterRegex.
func scanFilter(line string) []string {
	s := afterReplyCode(line)
	i := leadingSpace(s)
	if i < 0 {
		return nil
	}

	var direction string
	switch {
	case strings.HasPrefix(s[i:], filterInput+filterKey):
		direction = filterInput
	case strings.HasPrefix(s[i:], filterOutput+filterKey):
		direction = filterOutput
	default:
		return nil
	}

	i += len(direction) + len(filterKey)
	value := skipSpace(s, i)
	if value == i {
		return nil
	}

	return []string{line, direction, s[value:]}
}

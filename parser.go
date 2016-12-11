/*
Copyright 2016 Daniel Czerwonk (d.czerwonk@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

var (
	protoRegex  *regexp.Regexp
	routeRegex  *regexp.Regexp
	uptimeRegex *regexp.Regexp
)

func initRegexes() {
	protoRegex, _ = regexp.Compile("^([^\\s]+)\\s+BGP\\s+([^\\s]+)\\s+([^\\s]+)\\s+([^\\s]+)\\s+(.*?)\\s*$")
	routeRegex, _ = regexp.Compile("^\\s+Routes:\\s+(\\d+) imported, \\d+ filtered, (\\d+) exported")
	uptimeRegex, _ = regexp.Compile("^(?:((\\d+):(\\d{2}):(\\d{2}))|\\d+)$")
}

func parseOutput(data []byte, ipVersion int) []*session {
	sessions := make([]*session, 0)

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	var current *session = nil

	for scanner.Scan() {
		line := scanner.Text()
		if session, ok := parseLineForSession(line, ipVersion); ok {
			current = session
			sessions = append(sessions, current)
		}

		if current != nil {
			parseLineForRoutes(line, current)
		}

		if line == "" {
			current = nil
		}
	}

	return sessions
}

func parseLineForSession(line string, ipVersion int) (*session, bool) {
	match := protoRegex.FindStringSubmatch(line)

	if match == nil {
		return nil, false
	}

	session := session{name: match[1], ipVersion: ipVersion, established: parseState(match[5]), uptime: parseUptime(match[4])}
	return &session, true
}

func parseLineForRoutes(line string, session *session) {
	match := routeRegex.FindStringSubmatch(line)

	if match != nil {
		session.imported, _ = strconv.ParseInt(match[1], 10, 64)
		session.exported, _ = strconv.ParseInt(match[2], 10, 64)
	}
}

func parseState(state string) int {
	if state == "Established" {
		return 1
	} else {
		return 0
	}
}

func parseUptime(value string) int {
	match := uptimeRegex.FindStringSubmatch(value)

	if match == nil {
		return 0
	}

	if match[1] != "" {
		return parseUptimeForDuration(match)
	}

	return parseUptimeForTimestamp(value)
}

func parseUptimeForDuration(duration []string) int {
	h := parseInt(duration[2])
	m := parseInt(duration[3])
	s := parseInt(duration[4])
	str := fmt.Sprintf("%dh%dm%ds", h, m, s)

	d, err := time.ParseDuration(str)

	if err != nil {
		log.Println(err)
		return 0
	}

	return int(d.Seconds())
}

func parseUptimeForTimestamp(timestamp string) int {
	since := parseInt(timestamp)

	s := time.Unix(since, 0)
	d := time.Since(s)
	return int(d.Seconds())
}

func parseInt(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		log.Println(err)
		return 0
	}

	return i
}

package parser

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/stretchr/testify/require"
)

var scanSeeds = []string{
	"1002-peer BGP master up 2026-10-07 13:24:30.123 Established\n" +
		"1006-  Description:    site=ams role=upstream\n" +
		"  Channel ipv4\n" +
		"    State:          UP\n" +
		"    Input filter:   ACCEPT\n" +
		"    Output filter:  REJECT\n" +
		"    Routes:         1 imported, 2 filtered, 3 exported, 1 preferred\n" +
		"    Route change stats:     received   rejected   filtered    ignored   accepted\n" +
		"      Import updates:              1          2          3          4          5\n" +
		"      Import withdraws:            6          7          8          9        ---\n" +
		"      Export updates:             11         12         13         14         15         16         17\n" +
		"      Export withdraws:           16         17         18         19        ---\n" +
		"  Channel ipv6\n" +
		"    Routes:         3 imported, 4 exported\n" +
		"\n",
	"1002-st1 Static master4 up 2026-10-07 13:24:30\n" +
		"1006-  Channel ipv4\n" +
		"    Routes:         150 imported, 0 exported, 150 preferred\n",
	"ospf1 OSPF master up 00:01:00 Running\n  Channel roa6\n",
	"peer BGP master start 00:01:00 Active\n" +
		"1006-  Message:        Description: x Input filter: y Import updates: 9 9 9 9 9\n" +
		"  Message:\n",
	"kernel1 Kernel master up 1700000000\n",
	"1002- BGP t up x\n",
	"p BGP t up 2026-10-07 13:24:30x info\n",
	"p BGP t up 2026-10-07 13:24:30.12x  info  \n",
	"    Routes:\t5 imported, 6 exported, x preferred\n",
}

func FuzzParseProtocolsScanMatchesRegex(f *testing.F) {
	overrideNowFunc(func() time.Time {
		return time.Date(2026, 10, 7, 14, 0, 0, 0, time.Local)
	})

	for _, seed := range scanSeeds {
		f.Add(seed, "")
		f.Add(seed, "4")
	}
	if raw, err := os.ReadFile("testdata/protocols_all.raw"); err == nil {
		f.Add(string(raw), "")
	}

	f.Fuzz(func(t *testing.T, data string, ipVersion string) {
		want, wantErr := ParseProtocolsWithError([]byte(data), ipVersion)
		got, gotErr := ParseProtocolsScanWithError([]byte(data), ipVersion)

		require.Equal(t, wantErr, gotErr)
		require.Equal(t, want, got)
	})
}

func FuzzLineMatchersScanMatchRegex(f *testing.F) {
	for _, seed := range scanSeeds {
		for _, line := range strings.Split(seed, "\n") {
			f.Add(line)
		}
	}

	f.Fuzz(func(t *testing.T, line string) {
		require.Equal(t, regexMatchers.protocol(line), scanMatchers.protocol(line), "protocol")
		require.Equal(t, regexMatchers.description(line), scanMatchers.description(line), "description")
		require.Equal(t, regexMatchers.message(line), scanMatchers.message(line), "message")
		require.Equal(t, regexMatchers.channel(line), scanMatchers.channel(line), "channel")
		require.Equal(t, regexMatchers.routes(line), scanMatchers.routes(line), "routes")
		require.Equal(t, regexMatchers.routeChange(line), scanMatchers.routeChange(line), "routeChange")
		require.Equal(t, regexMatchers.filter(line), scanMatchers.filter(line), "filter")
	})
}

func BenchmarkParseProtocols(b *testing.B) {
	raw, err := os.ReadFile("testdata/protocols_all.raw")
	require.NoError(b, err)

	for _, bc := range []struct {
		name  string
		parse func([]byte, string) ([]*protocol.Protocol, error)
	}{
		{"regex", ParseProtocolsWithError},
		{"scan", ParseProtocolsScanWithError},
	} {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := bc.parse(raw, ""); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

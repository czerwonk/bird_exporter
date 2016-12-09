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
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type session struct {
	name        string
	ipVersion   int
	established int
	imported    int64
	exported    int64
	uptime      int
}

const version string = "0.2.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9200", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdClient    = flag.String("bird.client", "birdc", "Binary to communicate with the bird routing daemon")
)

func main() {
	initRegexes()
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("bird_bgp_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("BGP metric exporter for bird routing daemon")
}

func startServer() {
	fmt.Printf("Starting bird BGP exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	sessions := getSessions()

	for _, s := range sessions {
		writeLineForSession(s, w)
	}
}

func writeLineForSession(s *session, w io.Writer) {
	fmt.Fprintf(w, "bgp%d_session_up{name=\"%s\"} %d\n", s.ipVersion, s.name, s.established)
	fmt.Fprintf(w, "bgp%d_session_prefix_count_import{name=\"%s\"} %d\n", s.ipVersion, s.name, s.imported)
	fmt.Fprintf(w, "bgp%d_session_prefix_count_export{name=\"%s\"} %d\n", s.ipVersion, s.name, s.exported)
	fmt.Fprintf(w, "bgp%d_session_uptime{name=\"%s\"} %d\n", s.ipVersion, s.name, s.uptime)
}

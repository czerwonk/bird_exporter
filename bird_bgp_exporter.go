package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type session struct {
	name        string
	ipVersion   int
	established int
	imported    int64
	exported    int64
}

const version string = "0.1"

var (
	protoRegex    *regexp.Regexp
	routeRegex    *regexp.Regexp
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
	fmt.Println("BGP exporter for bird routing daemon")
}

func startServer() {
	fmt.Printf("Starting bird BGP exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func initRegexes() {
	protoRegex, _ = regexp.Compile("^([^\\s]+)\\s+BGP\\s+([^\\s]+)\\s+([^\\s]+)\\s+([\\d]+)\\s+(.*?)\\s*$")
	routeRegex, _ = regexp.Compile("^\\s+Routes:\\s+(\\d+) imported, \\d+ filtered, (\\d+) exported")
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
}

func getSessions() []*session {
	birdSessions := getSessionsFromBird(4)
	bird6Sessions := getSessionsFromBird(6)

	return append(birdSessions, bird6Sessions...)
}

func getSessionsFromBird(ipVersion int) []*session {
	client := *birdClient

	if ipVersion == 6 {
		client += "6"
	}

	output := getBirdOutput(client)
	return parseOutput(output, ipVersion)
}

func getBirdOutput(birdClient string) []byte {
	b, err := exec.Command(birdClient, "show", "protocols", "all").Output()

	if err != nil {
		b = make([]byte, 0)
		log.Println(err)
	}

	return b
}

func parseOutput(data []byte, ipVersion int) []*session {
	sessions := make([]*session, 0)

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	var current *session = nil

	for scanner.Scan() {
		line := scanner.Text()
		if session, res := parseLineForSession(line, ipVersion); res == true {
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

	if match != nil {
		session := session{name: match[1], ipVersion: ipVersion, established: parseState(match[5])}
		return &session, true
	}

	return nil, false
}

func parseLineForRoutes(line string, session *session) {
	match := routeRegex.FindStringSubmatch(line)

	if match != nil {
		session.imported, _ = strconv.ParseInt(match[1], 0, 64)
		session.exported, _ = strconv.ParseInt(match[2], 0, 64)
	}
}

func parseState(state string) int {
	if state == "Established" {
		return 1
	} else {
		return 0
	}
}

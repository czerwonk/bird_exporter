package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

type session struct {
	name        string
	ipVersion   int
	established int
}

var regex *regexp.Regexp

func main() {
	r, _ := regexp.Compile("^([^\\s]+)[\\s]+BGP[\\s]+([^\\s]+)[\\s]+([^\\s]+)[\\s]+([\\d]+)[\\s]+(.*?)[\\s]*$")
	regex = r

	fmt.Println("Starting bgp exporter")
	http.HandleFunc("/metrics", handleMetricsRequest)

	log.Fatal(http.ListenAndServe(":9200", nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	sessions := getSessions()

	for _, s := range sessions {
		writeLineForSession(s, w)
	}
}

func writeLineForSession(s session, w io.Writer) {
	fmt.Fprintf(w, "bgp%d_session_up{name=\"%s\"} %d\n", s.ipVersion, s.name, s.established)
}

func getSessions() []session {
	birdSessions := getSessionsFromBird(4)
	bird6Sessions := getSessionsFromBird(6)

	return append(birdSessions, bird6Sessions...)
}

func getSessionsFromBird(ipVersion int) []session {
	client := "birdc"

	if ipVersion == 6 {
		client += "6"
	}

	output := getBirdOutput(client)
	return parseOutput(output, ipVersion)
}

func getBirdOutput(birdClient string) []byte {
	b, err := exec.Command(birdClient, "show", "protocols").Output()

	if err != nil {
		b = make([]byte, 0)
		log.Print(err)
	}

	return b
}

func parseOutput(data []byte, ipVersion int) []session {
	sessions := make([]session, 0)

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if session, res := parseLine(line, ipVersion); res == true {
			sessions = append(sessions, *session)
		}
	}

	return sessions
}

func parseLine(line string, ipVersion int) (*session, bool) {
	match := regex.FindStringSubmatch(line)

	if match != nil {
		session := session{name: match[1], ipVersion: ipVersion, established: parseState(match[4])}
		return &session, true
	} else {
		return nil, false
	}
}

func parseState(state string) int {
	if state == "Established" {
		return 1
	} else {
		return 0
	}
}

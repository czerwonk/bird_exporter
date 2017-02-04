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
)

const version string = "0.5"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9200", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdClient    = flag.String("bird.client", "birdc", "Binary to communicate with the bird routing daemon")
	birdEnabled   = flag.Bool("bird.ipv4", true, "Get protocols from bird")
	bird6Enabled  = flag.Bool("bird.ipv6", true, "Get protocols from bird6")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("bird_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Metric exporter for bird routing daemon")
}

func startServer() {
	fmt.Printf("Starting bird exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(io.Writer, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)
		err := f(wr, r)
		wr.Flush()

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(buf.Bytes())

		if err != nil {
			log.Println(err)
		}
	}
}

func handleMetricsRequest(w io.Writer, r *http.Request) error {
	protocols, err := getProtocols()
	if err != nil {
		return err
	}

	for _, s := range protocols {
		writeForBgpSession(s, w)
	}

	return nil
}

func writeForBgpSession(s *protocol, w io.Writer) {
	fmt.Fprintf(w, "bgp%d_session_up{name=\"%s\"} %d\n", s.ipVersion, s.name, s.established)
	fmt.Fprintf(w, "bgp%d_session_prefix_count_import{name=\"%s\"} %d\n", s.ipVersion, s.name, s.imported)
	fmt.Fprintf(w, "bgp%d_session_prefix_count_export{name=\"%s\"} %d\n", s.ipVersion, s.name, s.exported)
	fmt.Fprintf(w, "bgp%d_session_uptime{name=\"%s\"} %d\n", s.ipVersion, s.name, s.uptime)
}

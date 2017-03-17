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

const version string = "0.5.2"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9200", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdClient    = flag.String("bird.client", "birdc", "Binary to communicate with the bird routing daemon")
	birdEnabled   = flag.Bool("bird.ipv4", true, "Get protocols from bird")
	bird6Enabled  = flag.Bool("bird.ipv6", true, "Get protocols from bird6")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: bird_exporter [ ... ]\n\nParameters:\n")
		flag.PrintDefaults()
	}
}

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

	for _, p := range protocols {
		switch p.proto {
		case BGP:
			writeForBgpSession(p, w)
		case OSPF:
			writeForOspf(p, w)
		}
	}

	return nil
}

func writeForBgpSession(p *protocol, w io.Writer) {
	prefix := fmt.Sprintf("bgp%d_session", p.ipVersion)
	writeForProtocol(p, prefix, w)
}

func writeForOspf(p *protocol, w io.Writer) {
	if p.ipVersion == 4 {
		writeForProtocol(p, "ospf", w)
	} else {
		writeForProtocol(p, "ospfv3", w)
	}
}

func writeForProtocol(p *protocol, prefix string, w io.Writer) {
	fmt.Fprintf(w, "%s_up{name=\"%s\"} %d\n", prefix, p.name, p.up)
	fmt.Fprintf(w, "%s_prefix_count_import{name=\"%s\"} %d\n", prefix, p.name, p.imported)
	fmt.Fprintf(w, "%s_prefix_count_export{name=\"%s\"} %d\n", prefix, p.name, p.exported)
	fmt.Fprintf(w, "%s_uptime{name=\"%s\"} %d\n", prefix, p.name, p.uptime)

	for k, v := range p.attributes {
		fmt.Fprintf(w, "%s_%s{name=\"%s\"} %v\n", prefix, k, v)
	}
}

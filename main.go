package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "0.8.2"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9324", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdSocket    = flag.String("bird.socket", "/var/run/bird.ctl", "Socket to communicate with bird routing daemon")
	bird6Socket   = flag.String("bird.socket6", "/var/run/bird6.ctl", "Socket to communicate with bird6 routing daemon")
	birdEnabled   = flag.Bool("bird.ipv4", true, "Get protocols from bird")
	bird6Enabled  = flag.Bool("bird.ipv6", true, "Get protocols from bird6")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: bird_exporter [ ... ]\n\nParameters:")
		fmt.Println()
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
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Metric exporter for bird routing daemon")
}

func startServer() {
	log.Infof("Starting bird exporter (Version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Bird Routing Daemon Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>Bird Routing Daemon Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/czerwonk/bird_exporter">github.com/czerwonk/bird_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	protocols, err := getProtocols()
	if err != nil {
		return err
	}

	if len(protocols) > 0 {
		reg := prometheus.NewRegistry()
		c := NewMetricCollectorForProtocols(protocols)
		reg.MustRegister(c)

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
	}

	return nil
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "1.2.3"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9324", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdSocket    = flag.String("bird.socket", "/var/run/bird.ctl", "Socket to communicate with bird routing daemon")
	birdV2        = flag.Bool("bird.v2", false, "Bird major version >= 2.0 (multi channel protocols)")
	newFormat     = flag.Bool("format.new", false, "New metric format (more convenient / generic)")
	enableBgp     = flag.Bool("proto.bgp", true, "Enables metrics for protocol BGP")
	enableOspf    = flag.Bool("proto.ospf", true, "Enables metrics for protocol OSPF")
	enableKernel  = flag.Bool("proto.kernel", true, "Enables metrics for protocol Kernel")
	enableStatic  = flag.Bool("proto.static", true, "Enables metrics for protocol Static")
	enableDirect  = flag.Bool("proto.direct", true, "Enables metrics for protocol Direct")
	// pre bird 2.0
	bird6Socket  = flag.String("bird.socket6", "/var/run/bird6.ctl", "Socket to communicate with bird6 routing daemon (not compatible with -bird.v2)")
	birdEnabled  = flag.Bool("bird.ipv4", true, "Get protocols from bird (not compatible with -bird.v2)")
	bird6Enabled = flag.Bool("bird.ipv6", true, "Get protocols from bird6 (not compatible with -bird.v2)")
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

	if !*newFormat {
		log.Info("INFO: You are using the old metric format. Please consider using the new (more convenient one) by setting -format.new=true.")
	}

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
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()
	p := enabledProtocols()
	c := NewMetricCollector(*newFormat, p)
	reg.MustRegister(c)

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
func enabledProtocols() int {
	res := 0

	if *enableBgp {
		res |= protocol.BGP
	}
	if *enableOspf {
		res |= protocol.OSPF
	}
	if *enableKernel {
		res |= protocol.Kernel
	}
	if *enableStatic {
		res |= protocol.Static
	}
	if *enableDirect {
		res |= protocol.Direct
	}

	return res
}

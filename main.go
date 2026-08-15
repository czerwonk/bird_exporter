package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var version = "1.5.0"

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverReadTimeout       = 10 * time.Second
	serverWriteTimeout      = 30 * time.Second
	serverIdleTimeout       = 60 * time.Second
	serverMaxHeaderBytes    = 1 << 20
)

var (
	showVersion      = flag.Bool("version", false, "Print version information.")
	listenAddress    = flag.String("web.listen-address", ":9324", "Address on which to expose metrics and web interface.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	birdSocket       = flag.String("bird.socket", "/var/run/bird.ctl", "Socket to communicate with bird routing daemon")
	birdV2           = flag.Bool("bird.v2", false, "Bird major version >= 2.0 (multi channel protocols)")
	tlsEnabled       = flag.Bool("tls.enabled", false, "Enables TLS")
	tlsCertChainPath = flag.String("tls.cert-file", "", "Path to TLS cert file")
	tlsKeyPath       = flag.String("tls.key-file", "", "Path to TLS key file")
	newFormat        = flag.Bool("format.new", true, "New metric format (more convenient / generic)")
	enableBGP        = flag.Bool("proto.bgp", true, "Enables metrics for protocol BGP")
	enableOSPF       = flag.Bool("proto.ospf", true, "Enables metrics for protocol OSPF")
	enableKernel     = flag.Bool("proto.kernel", true, "Enables metrics for protocol Kernel")
	enableStatic     = flag.Bool("proto.static", true, "Enables metrics for protocol Static")
	enableDirect     = flag.Bool("proto.direct", true, "Enables metrics for protocol Direct")
	enableBabel      = flag.Bool("proto.babel", true, "Enables metrics for protocol Babel")
	enableRPKI       = flag.Bool("proto.rpki", true, "Enables metrics for protocol RPKI")
	enableBFD        = flag.Bool("proto.bfd", true, "Enables metrics for protocol BFD")
	// pre bird 2.0
	bird6Socket            = flag.String("bird.socket6", "/var/run/bird6.ctl", "Socket to communicate with bird6 routing daemon (not compatible with -bird.v2)")
	birdEnabled            = flag.Bool("bird.ipv4", true, "Get protocols from bird (not compatible with -bird.v2)")
	bird6Enabled           = flag.Bool("bird.ipv6", true, "Get protocols from bird6 (not compatible with -bird.v2)")
	descriptionLabels      = flag.Bool("format.description-labels", false, "Add labels from protocol descriptions.")
	descriptionLabelsRegex = flag.String("format.description-labels-regex", "(\\w+)=(\\w+)", "Regex to extract labels from protocol description")
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
	log.Infof("Starting bird exporter (Version: %s)", version)
	if err := validateListenAddress(*listenAddress); err != nil {
		log.Fatal(err)
	}

	if !*newFormat {
		log.Info("INFO: You are using the old metric format. Please consider using the new (more convenient one) by setting -format.new=true.")
	}

	mux := newHTTPMux(*metricsPath, handleMetricsRequest)
	server := newHTTPServer(*listenAddress, mux)

	log.Infof("Listening for %s on %s (TLS: %v)", *metricsPath, *listenAddress, *tlsEnabled)
	if *tlsEnabled {
		log.Fatal(server.ListenAndServeTLS(*tlsCertChainPath, *tlsKeyPath))
		return
	}

	log.Fatal(server.ListenAndServe())
}

func newHTTPMux(metricsPath string, metricsHandler http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(w, request)
			return
		}
		if request.Method != http.MethodGet && request.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := io.WriteString(w, `<html>
			<head><title>Bird Routing Daemon Exporter (Version `+version+`)</title></head>
			<body>
			<h1>Bird Routing Daemon Exporter</h1>
			<p><a href="`+metricsPath+`">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/czerwonk/bird_exporter">github.com/czerwonk/bird_exporter</a></p>
			</body>
			</html>`)
		if err != nil {
			log.Warnf("Could not write landing page: %v", err)
		}
	})
	mux.HandleFunc(metricsPath, func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		metricsHandler(w, request)
	})
	return mux
}

func newHTTPServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
		MaxHeaderBytes:    serverMaxHeaderBytes,
	}
}

func validateListenAddress(address string) error {
	if address == "" || strings.TrimSpace(address) != address {
		return fmt.Errorf("web listen address must be a non-empty host:port: %q", address)
	}

	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid web listen address %q: %w", address, err)
	}
	if port == "" {
		return fmt.Errorf("web listen address must include a port: %q", address)
	}
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("invalid web listen port %q: %w", port, err)
	}

	return nil
}

func newPrometheusHandler(gatherer prometheus.Gatherer) http.Handler {
	l := log.New()
	l.Level = log.ErrorLevel

	return promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{
		ErrorLog:      l,
		ErrorHandling: promhttp.HTTPErrorOnError,
	})
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()
	p := enabledProtocols()
	c := NewMetricCollector(*newFormat, p, *descriptionLabels, *birdSocket)
	reg.MustRegister(c)

	newPrometheusHandler(reg).ServeHTTP(w, r)
}

func enabledProtocols() protocol.Proto {
	res := protocol.Proto(0)

	if *enableBGP {
		res |= protocol.BGP
	}
	if *enableOSPF {
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
	if *enableBabel {
		res |= protocol.Babel
	}
	if *enableRPKI {
		res |= protocol.RPKI
	}
	if *enableBFD {
		res |= protocol.BFD
	}

	return res
}

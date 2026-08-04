package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

var invalidMetricDesc = prometheus.NewDesc(
	"bird_exporter_invalid_test_metric",
	"Test-only invalid metric.",
	nil,
	nil,
)

type invalidCollector struct{}

func (*invalidCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- invalidMetricDesc
}

func (*invalidCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.NewInvalidMetric(invalidMetricDesc, errors.New("invalid test metric"))
}

func TestHTTPServerHasResourceTimeouts(t *testing.T) {
	server := newHTTPServer("127.0.0.1:0", http.NewServeMux())

	if server.ReadHeaderTimeout != serverReadHeaderTimeout ||
		server.ReadTimeout != serverReadTimeout ||
		server.WriteTimeout != serverWriteTimeout ||
		server.IdleTimeout != serverIdleTimeout ||
		server.MaxHeaderBytes != serverMaxHeaderBytes {
		t.Fatalf("unexpected server limits: %+v", server)
	}
}

func TestHTTPMuxRejectsUnexpectedRoutesAndMethods(t *testing.T) {
	metricsCalled := false
	mux := newHTTPMux("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		metricsCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantAllow  string
	}{
		{name: "root", method: http.MethodGet, path: "/", wantStatus: http.StatusOK},
		{name: "unknown path", method: http.MethodGet, path: "/debug", wantStatus: http.StatusNotFound},
		{name: "root method", method: http.MethodPost, path: "/", wantStatus: http.StatusMethodNotAllowed, wantAllow: "GET, HEAD"},
		{name: "metrics method", method: http.MethodPost, path: "/metrics", wantStatus: http.StatusMethodNotAllowed, wantAllow: "GET"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.path, nil)
			mux.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if response.Header().Get("Allow") != test.wantAllow {
				t.Fatalf("Allow = %q, want %q", response.Header().Get("Allow"), test.wantAllow)
			}
		})
	}

	if metricsCalled {
		t.Fatal("metrics handler called for rejected method")
	}
}

func TestHTTPMuxServesMetrics(t *testing.T) {
	metricsCalled := false
	mux := newHTTPMux("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		metricsCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	mux.ServeHTTP(response, request)

	if !metricsCalled || response.Code != http.StatusNoContent {
		t.Fatalf("metricsCalled = %v, status = %d", metricsCalled, response.Code)
	}
}

func TestPrometheusHandlerFailsClosedOnGatherError(t *testing.T) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(&invalidCollector{})

	response := httptest.NewRecorder()
	newPrometheusHandler(registry).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/metrics", nil),
	)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

func TestValidateListenAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "all interfaces", address: ":9324"},
		{name: "IPv4 loopback", address: "127.0.0.1:9324"},
		{name: "IPv6 loopback", address: "[::1]:9324"},
		{name: "service port", address: "localhost:http"},
		{name: "empty", address: "", wantErr: true},
		{name: "missing port", address: "127.0.0.1", wantErr: true},
		{name: "empty port", address: "127.0.0.1:", wantErr: true},
		{name: "unbracketed IPv6", address: "::1:9324", wantErr: true},
		{name: "unknown service", address: "127.0.0.1:not-a-service", wantErr: true},
		{name: "surrounding whitespace", address: " 127.0.0.1:9324", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateListenAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateListenAddress(%q) error = %v, wantErr = %v", tt.address, err, tt.wantErr)
			}
		})
	}
}

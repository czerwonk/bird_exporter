package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

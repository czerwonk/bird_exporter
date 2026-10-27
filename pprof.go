package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	log "github.com/sirupsen/logrus"
)

// Temporary profiling hook: set BIRD_EXPORTER_PPROF=127.0.0.1:6060 to serve
// /debug/pprof/ on a separate listener.
func init() {
	addr := os.Getenv("BIRD_EXPORTER_PPROF")
	if addr == "" {
		return
	}

	go func() {
		log.Infof("pprof listening on %s", addr)
		log.Error(http.ListenAndServe(addr, nil))
	}()
}

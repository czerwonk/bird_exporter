package main

import (
	"github.com/KimMachineGun/automemlimit/memlimit"
	log "github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, err := maxprocs.Set(maxprocs.Logger(log.Debugf))
	if err != nil {
		log.Error("failed to setup CPU limits, continue without limits: %w", err)
	}
	_, err = memlimit.Set(
		memlimit.WithRatio(0.9),
		memlimit.WithProvider(
			memlimit.ApplyFallback(
				memlimit.FromCgroup,
				memlimit.FromSystem,
			),
		),
	)
	if err != nil {
		log.Error("failed to setup memory limits, continue without limits: %w", err)
	}
}

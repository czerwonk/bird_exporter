package main

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// cachingGatherer queries bird at most once per ttl. All requests within the
// ttl are served from the same snapshot, stamped with the time it was taken.
// Failed snapshots are cached as well, so an unreachable bird is not hammered.
type cachingGatherer struct {
	gatherer prometheus.Gatherer
	ttl      time.Duration
	now      func() time.Time

	mu       sync.Mutex
	families []*dto.MetricFamily
	err      error
	expires  time.Time
}

func newCachingGatherer(gatherer prometheus.Gatherer, ttl time.Duration) *cachingGatherer {
	return &cachingGatherer{
		gatherer: gatherer,
		ttl:      ttl,
		now:      time.Now,
	}
}

// Gather holds the lock while refreshing, so concurrent requests wait for the
// running refresh instead of starting their own.
func (c *cachingGatherer) Gather() ([]*dto.MetricFamily, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.now().Before(c.expires) {
		return c.families, c.err
	}

	takenAt := c.now()
	families, err := c.gatherer.Gather()
	stampFamilies(families, takenAt)

	c.families, c.err = families, err
	c.expires = takenAt.Add(c.ttl)

	return c.families, c.err
}

func stampFamilies(families []*dto.MetricFamily, t time.Time) {
	ms := t.UnixMilli()
	for _, f := range families {
		for _, m := range f.GetMetric() {
			m.TimestampMs = &ms
		}
	}
}

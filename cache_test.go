package main

import (
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestCachingGathererServesSnapshotWithinTTL(t *testing.T) {
	calls := 0
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "test_value", Help: "Test value."}, func() float64 {
		calls++
		return float64(calls)
	}))

	now := time.Unix(1000, 0)
	c := newCachingGatherer(reg, 10*time.Second)
	c.now = func() time.Time { return now }

	first := gatherValue(t, c)
	now = now.Add(9 * time.Second)
	second := gatherValue(t, c)
	if calls != 1 || first.GetGauge().GetValue() != 1 || second.GetGauge().GetValue() != 1 {
		t.Fatalf("expected one bird query within ttl, got %d", calls)
	}
	if second.GetTimestampMs() != time.Unix(1000, 0).UnixMilli() {
		t.Fatalf("expected snapshot timestamp, got %d", second.GetTimestampMs())
	}

	now = now.Add(time.Second)
	third := gatherValue(t, c)
	if calls != 2 || third.GetGauge().GetValue() != 2 {
		t.Fatalf("expected refresh after ttl, got %d queries", calls)
	}
	if third.GetTimestampMs() != now.UnixMilli() {
		t.Fatalf("expected refreshed timestamp, got %d", third.GetTimestampMs())
	}
}

func TestCachingGathererCachesFailures(t *testing.T) {
	calls := 0
	failing := prometheus.GathererFunc(func() ([]*dto.MetricFamily, error) {
		calls++
		return nil, errors.New("bird unreachable")
	})

	now := time.Unix(1000, 0)
	c := newCachingGatherer(failing, 10*time.Second)
	c.now = func() time.Time { return now }

	for range 3 {
		if _, err := c.Gather(); err == nil {
			t.Fatal("expected cached error")
		}
	}
	if calls != 1 {
		t.Fatalf("expected one query for cached failure, got %d", calls)
	}
}

func gatherValue(t *testing.T, g prometheus.Gatherer) *dto.Metric {
	t.Helper()

	families, err := g.Gather()
	if err != nil {
		t.Fatal(err)
	}
	if len(families) != 1 || len(families[0].GetMetric()) != 1 {
		t.Fatalf("unexpected families: %v", families)
	}

	return families[0].GetMetric()[0]
}

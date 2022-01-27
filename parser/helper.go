package parser

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	nowFunc func() time.Time
)

func overrideNowFunc(f func() time.Time) {
	nowFunc = f
}

func currentTime() time.Time {
	return nowFunc()
}

func parseInt(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		log.Errorln(err)
		return 0
	}

	return i
}

func parseFloat(value string) float64 {
	i, err := strconv.ParseFloat(value, 64)

	if err != nil {
		log.Errorln(err)
		return 0
	}

	return i
}

func parseUptimeForIso(s string) int {
	start, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		log.Errorln(err)
		return 0
	}

	return int(currentTime().Sub(start).Seconds())
}

func parseUptimeForDuration(duration []string) int {
	h := parseInt(duration[2])
	m := parseInt(duration[3])
	s := parseInt(duration[4])
	str := fmt.Sprintf("%dh%dm%ds", h, m, s)

	d, err := time.ParseDuration(str)
	if err != nil {
		log.Errorln(err)
		return 0
	}

	return int(d.Seconds())
}

func parseUptimeForTimestamp(timestamp string) int {
	since := parseInt(timestamp)

	s := time.Unix(since, 0)
	d := time.Since(s)
	return int(d.Seconds())
}

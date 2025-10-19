package parser

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/czerwonk/bird_exporter/routes"
)

func ParseExportedRoutes(data []byte) (rList []routes.Route, err error) {
	tablePrefix := regexp.MustCompile(`^(1007\-Table) (\w+?)\:`)
	routePrefix := regexp.MustCompile(`^ *(\d+\.\d+\.\d+\.\d+)\/(\d+) +(.+)`)
	routePostfix := regexp.MustCompile(`^(\w+) \[(.+?) +(.+?)\] +(\*? +)\((\d+)\)(?: \[(.+)\])?`)
	viaPrefix := regexp.MustCompile(`^ +\tvia (\d+\.\d+\.\d+\.\d+) on (\w+?)(?: weight (\d+))?$`)
	devPrefix := regexp.MustCompile(`^ +\tdev (\w+?)$`)

	var tName string
	rt := routes.Route{}
	rtt := routes.RouteTarget{}

	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " ")

		if m := tablePrefix.FindStringSubmatch(line); m != nil {
			tName = m[2]
			continue
		}

		if tName == "" {
			continue
		}

		if m := routePrefix.FindStringSubmatch(line); m != nil {
			if rtt.RouteType != "" {
				rt.Targets = append(rt.Targets, rtt)
			}

			if rt.Table != "" {
				rList = append(rList, rt)
				rt = routes.Route{}
			}

			rt.Table = tName
			rt.Prefix = m[1]
			rt.NetLen = m[2]

			if m2 := routePostfix.FindStringSubmatch(m[3]); m2 != nil {
				rtt.RouteType = m2[1]
				rtt.RouteSource = m2[2]
				rtt.FirstSeen = m2[3]
				if len(m2[4]) > 0 {
					rtt.IsBest = true
				}
				if v, e := strconv.ParseInt(m2[5], 10, 16); e == nil {
					rtt.Metric = int(v)
				}

				if len(m2) > 6 {
					// AS is specified
					rtt.LastAS = m2[6]
				}
			}
			continue
		}

		if m2 := routePostfix.FindStringSubmatch(line); m2 != nil {
			rtt.RouteType = m2[1]
			rtt.RouteSource = m2[2]
			rtt.FirstSeen = m2[3]
			if len(m2[4]) > 0 {
				rtt.IsBest = true
			}
			if v, e := strconv.ParseInt(m2[5], 10, 16); e == nil {
				rtt.Metric = int(v)
			}

			if len(m2) > 6 {
				// AS is specified
				rtt.LastAS = m2[6]
			}

			continue
		}

		if m := viaPrefix.FindStringSubmatch(line); m != nil {
			rtv := routes.RouteVia{
				Via: m[1],
				Dev: m[2],
			}
			if len(m) > 3 {
				if v, e := strconv.ParseInt(m[3], 10, 16); e == nil {
					rtv.Weight = int(v)
				}
			}

			// Process ECMP multiple via
			if rtt.RouteType == "" {
				// Check correct data
				if len(rt.Targets) > 0 && len(rt.Targets[len(rt.Targets)-1].Via) > 0 {
					rt.Targets[len(rt.Targets)-1].Via = append(rt.Targets[len(rt.Targets)-1].Via, rtv)
				}
			} else {
				rtt.Via = append(rtt.Via, rtv)
				rt.Targets = append(rt.Targets, rtt)
				rtt = routes.RouteTarget{}
			}
			continue
		}

		if m := devPrefix.FindStringSubmatch(line); m != nil {
			rtt.Dev = m[1]
			continue
		}

	}

	if rtt.RouteType != "" {
		rt.Targets = append(rt.Targets, rtt)
	}

	if rt.Table != "" {
		rList = append(rList, rt)
	}

	return
}

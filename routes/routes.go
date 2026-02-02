package routes

// RouteVia - implemented as slice for ECMP support
type RouteVia struct {
	Via    string
	Dev    string
	Weight int
}
type RouteTarget struct {
	RouteType   string
	RouteSource string
	FirstSeen   string
	IsBest      bool
	Metric      int
	LastAS      string
	Dev         string
	Via         []RouteVia
}

type Route struct {
	Table   string
	Prefix  string
	NetLen  string
	Targets []RouteTarget
}

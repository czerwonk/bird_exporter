package routes

type RouteTarget struct {
	RouteType   string
	RouteSource string
	FirstSeen   string
	IsBest      bool
	Metric      int
	LastAS      *string
	Dev         string
	Via         string
}

type Route struct {
	Table   string
	Prefix  string
	NetLen  string
	Targets []RouteTarget
}

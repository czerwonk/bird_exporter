package protocol

const (
	PROTO_UNKNOWN = 0
	BGP           = 1
	OSPF          = 2
)

type Protocol struct {
	Name            string
	IpVersion       int
	Proto           int
	Up              int
	Imported        int64
	Exported        int64
	Filtered        int64
	Preferred       int64
	Uptime          int
	Attributes      map[string]float64
	ImportUpdates   RouteChangeCount
	ImportWithdraws RouteChangeCount
	ExportUpdates   RouteChangeCount
	ExportWithdraws RouteChangeCount
}

type RouteChangeCount struct {
	Received int64
	Rejected int64
	Filtered int64
	Ignored  int64
	Accepted int64
}

func NewProtocol(name string, proto, ipVersion, uptime int) *Protocol {
	return &Protocol{Name: name, Proto: proto, IpVersion: ipVersion, Uptime: uptime, Attributes: make(map[string]float64)}
}

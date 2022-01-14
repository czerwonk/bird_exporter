package protocol

const (
	PROTO_UNKNOWN = Proto(0)
	BGP           = Proto(1)
	OSPF          = Proto(2)
	Kernel        = Proto(4)
	Static        = Proto(8)
	Direct        = Proto(16)
	Babel         = Proto(32)
	RPKI          = Proto(64)
)

type Proto int

type Protocol struct {
	Name            string
	Description     string
	IPVersion       string
	ImportFilter    string
	ExportFilter    string
	Proto           Proto
	Up              int
	State           string
	Imported        int64
	Exported        int64
	Filtered        int64
	Preferred       int64
	Uptime          int
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

func NewProtocol(name string, proto Proto, ipVersion string, uptime int) *Protocol {
	return &Protocol{Name: name, Proto: proto, IPVersion: ipVersion, Uptime: uptime}
}

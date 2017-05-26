package protocol

const (
	PROTO_UNKNOWN = 0
	BGP           = 1
	OSPF          = 2
)

type Protocol struct {
	Name       string
	IpVersion  int
	Proto      int
	Up         int
	Imported   int64
	Exported   int64
	Filtered   int64
	Uptime     int
	Attributes map[string]float64
}

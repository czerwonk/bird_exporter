package protocol

type BFDSession struct {
	ProtocolName string
	IP           string
	Interface    string
	Up           bool
	Since        int
	Interval     float64
	Timeout      float64
}

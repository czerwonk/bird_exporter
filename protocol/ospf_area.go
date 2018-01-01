package protocol

type OspfArea struct {
	Name                  string
	InterfaceCount        int64
	NeighborCount         int64
	NeighborAdjacentCount int64
}

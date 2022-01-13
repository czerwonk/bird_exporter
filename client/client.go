package client

import "github.com/czerwonk/bird_exporter/protocol"

// Client retrieves information from Bird routing daemon
type Client interface {

	// GetProtocols retrieves protocol information and statistics from bird
	GetProtocols() ([]*protocol.Protocol, error)

	// GetOSPFAreas retrieves OSPF specific information from bird
	GetOSPFAreas(protocol *protocol.Protocol) ([]*protocol.OspfArea, error)

	// GetBGPStates retrieves BGP state information from bird
	GetBGPStates(protocol *protocol.Protocol) (*protocol.BgpState, error)
}

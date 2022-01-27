package client

import "github.com/czerwonk/bird_exporter/protocol"

// Client retrieves information from Bird routing daemon
type Client interface {

	// GetProtocols retrieves protocol information and statistics from bird
	GetProtocols() ([]*protocol.Protocol, error)

	// GetOSPFAreas retrieves OSPF specific information from bird
	GetOSPFAreas(protocol *protocol.Protocol) ([]*protocol.OSPFArea, error)

	// GetBFDSessions retrieves BFD specific information from bird
	GetBFDSessions(protocol *protocol.Protocol) ([]*protocol.BFDSession, error)
}

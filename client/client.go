package client

import "github.com/czerwonk/bird_exporter/protocol"

type Client interface {

	// GetProtocols retrieves protocol information and statistics from bird
	GetProtocols() ([]*protocol.Protocol, error)

	// GetOspfArea retrieves OSPF specific information from bird
	GetOspfAreas(protocol *protocol.Protocol) ([]*protocol.OspfArea, error)
}

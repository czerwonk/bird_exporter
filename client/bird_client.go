package client

import (
	"fmt"

	"github.com/czerwonk/bird_exporter/parser"
	"github.com/czerwonk/bird_exporter/protocol"
	birdsocket "github.com/czerwonk/bird_socket"
)

// BirdClient communicates with the bird socket to retrieve information
type BirdClient struct {
	Options *BirdClientOptions
}

// BirdClientOptions defines options to connect to bird
type BirdClientOptions struct {
	BirdV2       bool
	BirdEnabled  bool
	Bird6Enabled bool
	BirdSocket   string
	Bird6Socket  string
}

// GetProtocols retrieves protocol information and statistics from bird
func (c *BirdClient) GetProtocols() ([]*protocol.Protocol, error) {
	ipVersions := make([]string, 0)
	if c.Options.BirdV2 {
		ipVersions = append(ipVersions, "")
	} else {
		if c.Options.BirdEnabled {
			ipVersions = append(ipVersions, "4")
		}

		if c.Options.Bird6Enabled {
			ipVersions = append(ipVersions, "6")
		}
	}

	return c.protocolsFromBird(ipVersions)
}

// GetOSPFAreas retrieves OSPF specific information from bird
func (c *BirdClient) GetOSPFAreas(protocol *protocol.Protocol) ([]*protocol.OspfArea, error) {
	sock := c.socketFor(protocol.IPVersion)
	b, err := birdsocket.Query(sock, fmt.Sprintf("show ospf %s", protocol.Name))
	if err != nil {
		return nil, err
	}

	return parser.ParseOspf(b), nil
}

// GetBGPStates retrieves BGP state information from bird
func (c *BirdClient) GetBGPStates(protocol *protocol.Protocol) (*protocol.BgpState, error) {
	sock := c.socketFor(protocol.IPVersion)
	b, err := birdsocket.Query(sock, fmt.Sprintf("show protocols all %s", protocol.Name))
	if err != nil {
		return nil, err
	}
	return parser.ParseBgpState(b), nil
}

func (c *BirdClient) protocolsFromBird(ipVersions []string) ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	for _, ipVersion := range ipVersions {
		sock := c.socketFor(ipVersion)
		s, err := c.protocolsFromSocket(sock, ipVersion)
		if err != nil {
			return nil, err
		}

		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func (c *BirdClient) protocolsFromSocket(socketPath string, ipVersion string) ([]*protocol.Protocol, error) {
	b, err := birdsocket.Query(socketPath, "show protocols all")
	if err != nil {
		return nil, err
	}

	return parser.ParseProtocols(b, ipVersion), nil
}

func (c *BirdClient) socketFor(ipVersion string) string {
	if !c.Options.BirdV2 && ipVersion == "6" {
		return c.Options.Bird6Socket
	}

	return c.Options.BirdSocket
}

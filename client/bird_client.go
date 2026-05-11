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
	afiFamilies := make([]string, 0)
	if c.Options.BirdV2 {
		afiFamilies = append(afiFamilies, "")
	} else {
		if c.Options.BirdEnabled {
			afiFamilies = append(afiFamilies, "4")
		}

		if c.Options.Bird6Enabled {
			afiFamilies = append(afiFamilies, "6")
		}
	}

	return c.protocolsFromBird(afiFamilies)
}

// GetOSPFAreas retrieves OSPF specific information from bird
func (c *BirdClient) GetOSPFAreas(protocol *protocol.Protocol) ([]*protocol.OSPFArea, error) {
	sock := c.socketFor(protocol.AFIFamily)
	b, err := birdsocket.Query(sock, fmt.Sprintf("show ospf %s", protocol.Name))
	if err != nil {
		return nil, err
	}

	return parser.ParseOSPF(b), nil
}

// GetBFDSessions retrieves BFD specific information from bird
func (c *BirdClient) GetBFDSessions(protocol *protocol.Protocol) ([]*protocol.BFDSession, error) {
	sock := c.socketFor(protocol.AFIFamily)
	b, err := birdsocket.Query(sock, fmt.Sprintf("show bfd sessions %s", protocol.Name))
	if err != nil {
		return nil, err
	}

	return parser.ParseBFDSessions(protocol.Name, b), nil
}

func (c *BirdClient) protocolsFromBird(afiFamilies []string) ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	for _, afiFamily := range afiFamilies {
		sock := c.socketFor(afiFamily)
		s, err := c.protocolsFromSocket(sock, afiFamily)
		if err != nil {
			return nil, err
		}

		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func (c *BirdClient) protocolsFromSocket(socketPath string, afiFamily string) ([]*protocol.Protocol, error) {
	b, err := birdsocket.Query(socketPath, "show protocols all")
	if err != nil {
		return nil, err
	}

	return parser.ParseProtocols(b, afiFamily), nil
}

func (c *BirdClient) socketFor(afiFamily string) string {
	if !c.Options.BirdV2 && afiFamily == "6" {
		return c.Options.Bird6Socket
	}

	return c.Options.BirdSocket
}

// StatusFromSocket retrieves status information from bird
func (c *BirdClient) StatusFromSocket(socketPath string) (*parser.Status, error) {
	b, err := birdsocket.Query(socketPath, "show status")
	if err != nil {
		return nil, err
	}

	return parser.ParseStatus(b), nil
}

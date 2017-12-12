package main

import (
	"github.com/czerwonk/bird_exporter/parser"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/czerwonk/bird_socket"
)

func getProtocols() ([]*protocol.Protocol, error) {
	var protocols []*protocol.Protocol = nil
	var err error = nil

	if *birdV2 {
		protocols, err = getProtocolsFromBird(*birdSocket, "")
	} else {
		protocols, err = getProtocolsFromBird1()
	}

	return protocols, err
}

func getProtocolsFromBird1() ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	if *birdEnabled {
		s, err := getProtocolsFromBird(*birdSocket, "4")
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	if *bird6Enabled {
		s, err := getProtocolsFromBird(*bird6Socket, "6")
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func getProtocolsFromBird(socketPath string, ipVersion string) ([]*protocol.Protocol, error) {
	b, err := birdsocket.Query(socketPath, "show protocols all")
	if err != nil {
		return nil, err
	}

	return parser.Parse(b, ipVersion), nil
}

package main

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/czerwonk/bird_socket"
)

func getProtocols() ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	if *birdEnabled {
		s, err := getProtocolsFromBird(*birdSocket, 4)
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	if *bird6Enabled {
		s, err := getProtocolsFromBird(*bird6Socket, 6)
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func getProtocolsFromBird(socketPath string, ipVersion int) ([]*protocol.Protocol, error) {
	b, err := birdsocket.Query(socketPath, "show protocols all")
	if err != nil {
		return nil, err
	}

	return parseOutput(b, ipVersion), nil
}

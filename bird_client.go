package main

import (
	"os/exec"

	"github.com/czerwonk/bird_exporter/protocol"
)

func getProtocols() ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	if *birdEnabled {
		s, err := getProtocolsFromBird(4)
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	if *bird6Enabled {
		s, err := getProtocolsFromBird(6)
		if err != nil {
			return nil, err
		}
		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func getProtocolsFromBird(ipVersion int) ([]*protocol.Protocol, error) {
	client := *birdClient

	if ipVersion == 6 {
		client += "6"
	}

	output, err := getBirdOutput(client)
	if err != nil {
		return nil, err
	}

	return parseOutput(output, ipVersion), nil
}

func getBirdOutput(birdClient string) ([]byte, error) {
	return exec.Command(birdClient, "show", "protocols", "all").Output()
}

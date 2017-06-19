package main

import (
	"net"

	"regexp"

	"time"

	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/common/log"
)

const bufferSize = 4096
const timeout = 30

var birdReturnCodeRegex *regexp.Regexp

func init() {
	birdReturnCodeRegex = regexp.MustCompile("\\d{4} \n$")
}

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
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))

	buf := make([]byte, bufferSize)
	n, err := conn.Read(buf[:])
	if err != nil {
		return nil, err
	}
	log.Debug(string(buf[:n]))

	_, err = conn.Write([]byte("show protocols all\n"))
	if err != nil {
		return nil, err
	}

	output, err := readFromSocket(conn)
	if err != nil {
		return nil, err
	}

	return parseOutput(output, ipVersion), nil
}

func readFromSocket(conn net.Conn) ([]byte, error) {
	b := make([]byte, 0)
	buf := make([]byte, bufferSize)

	done := false
	for !done {
		n, err := conn.Read(buf[:])
		if err != nil {
			return nil, err
		}

		b = append(b, buf[:n]...)
		done = endsWithBirdReturnCode(b)
	}

	return b, nil
}

func endsWithBirdReturnCode(b []byte) bool {
	if len(b) < 6 {
		return false
	}

	return birdReturnCodeRegex.Match(b[len(b)-6:])
}

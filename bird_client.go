package main

import "os/exec"

func getSessions() ([]*session, error) {
	sessions := make([]*session, 0)

	if *birdEnabled {
		s, err := getSessionsFromBird(4)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s...)
	}

	if *bird6Enabled {
		s, err := getSessionsFromBird(6)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s...)
	}

	return sessions, nil
}

func getSessionsFromBird(ipVersion int) ([]*session, error) {
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

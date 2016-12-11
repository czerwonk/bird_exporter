/*
Copyright 2016 Daniel Czerwonk (d.czerwonk@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

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

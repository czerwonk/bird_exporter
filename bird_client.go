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

import (
	"log"
	"os/exec"
)

func getSessions() []*session {
	birdSessions := getSessionsFromBird(4)
	bird6Sessions := getSessionsFromBird(6)

	return append(birdSessions, bird6Sessions...)
}

func getSessionsFromBird(ipVersion int) []*session {
	client := *birdClient

	if ipVersion == 6 {
		client += "6"
	}

	output := getBirdOutput(client)
	return parseOutput(output, ipVersion)
}

func getBirdOutput(birdClient string) []byte {
	b, err := exec.Command(birdClient, "show", "protocols", "all").Output()

	if err != nil {
		b = make([]byte, 0)
		log.Println(err)
	}

	return b
}

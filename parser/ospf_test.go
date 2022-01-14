package parser

import (
	"testing"

	"github.com/czerwonk/testutils/assert"
)

func TestOSPFArea(t *testing.T) {
	data := "ospf1:\n" +
		"RFC1583 compatibility: disabled\n" +
		"Stub router: No\n" +
		"RT scheduler tick: 1\n" +
		"Number of areas: 2\n" +
		"Number of LSAs in DB:   33\n" +
		"    Area: 0.0.0.0 (0) [BACKBONE]\n" +
		"        Stub:   No\n" +
		"        NSSA:   No\n" +
		"        Transit:    No\n" +
		"        Number of interfaces:   3\n" +
		"        Number of neighbors:    2\n" +
		"        Number of adjacent neighbors:   1\n" +
		"    Area: 0.0.0.1 (1)\n" +
		"        Stub:   No\n" +
		"        NSSA:   No\n" +
		"        Transit:    No\n" +
		"        Number of interfaces:   4\n" +
		"        Number of neighbors:    6\n" +
		"        Number of adjacent neighbors:   5\n"
	a := ParseOSPF([]byte(data))
	assert.IntEqual("areas", 2, len(a), t)

	a1 := a[0]
	assert.StringEqual("Area1 Name", "0", a1.Name, t)
	assert.Int64Equal("Area1 InterfaceCount", 3, a1.InterfaceCount, t)
	assert.Int64Equal("Area1 NeighborCount", 2, a1.NeighborCount, t)
	assert.Int64Equal("Area1 NeighborAdjacentCount", 1, a1.NeighborAdjacentCount, t)

	a2 := a[1]
	assert.StringEqual("Area2 Name", "1", a2.Name, t)
	assert.Int64Equal("Area2 InterfaceCount", 4, a2.InterfaceCount, t)
	assert.Int64Equal("Area2 NeighborCount", 6, a2.NeighborCount, t)
	assert.Int64Equal("Area2 NeighborAdjacentCount", 5, a2.NeighborAdjacentCount, t)
}

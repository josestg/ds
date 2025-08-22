package ds

import (
	"math/rand"
	"testing"

	"github.com/josestg/ds/adt/adttest"
)

func TestHashMap(t *testing.T) {
	c := NewHashMap[int, int]
	kg := func() int {
		return rand.Intn(128)
	}
	vg := func() int {
		return rand.Intn(128)
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "hashmap", simulator: adttest.HashMapSimulator(c, kg, vg)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}

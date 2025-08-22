package ds

import (
	"math/rand"
	"testing"

	"github.com/josestg/ds/adt/adttest"
)

func TestBinarySearchTree(t *testing.T) {
	c := NewBinarySearchTree[int]
	g := func() int {
		return rand.Intn(128)
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "insert", simulator: adttest.BSTInsertSimulator(c, g)},
		{name: "in order", simulator: adttest.BSTInOrderSimulator(c, g)},
		{name: "min max", simulator: adttest.BSTMinMaxSimulator(c, g)},
		{name: "delete", simulator: adttest.BSTDeleteSimulator(c, g)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}

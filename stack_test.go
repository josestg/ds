package ds

import (
	"math/rand"
	"testing"

	"github.com/josestg/ds/adt/adttest"
)

func TestStack(t *testing.T) {
	c := NewStack[int]
	g := func() int {
		return rand.Intn(128)
	}

	simulator := adttest.StackSimulator(c, g)
	simulator.Run(t)
}

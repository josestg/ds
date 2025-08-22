package ds

import (
	"math/rand"
	"testing"

	"github.com/josestg/ds/adt/adttest"
)

func TestQueue(t *testing.T) {
	c := NewQueue[int]
	g := func() int {
		return rand.Intn(128)
	}
	simulator := adttest.QueueSimulator(c, g)
	simulator.Run(t)
}

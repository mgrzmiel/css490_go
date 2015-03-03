// CSS 490
// Magdalena Grzmiel
// Assignments #5
// This is a unit test for counter package

package counter

import (
	"sync"
	"testing"
)

const (
	// Symbolic constants to let the compiler find typos.
	Zeus   = "zeus"
	Hera   = "hera"
	Ares   = "ares"
	Athena = "athena"
)

func body(c *Counter, t *testing.T) {
	// Spawn off 4 concurrent threads and wait until they
	// complete.
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		c.Incr(Zeus, 2)
		c.Incr(Hera, 1)
		c.Incr(Athena, 1)
		c.Incr(Ares, 1)
		wg.Done()
	}()
	go func() {
		c.Incr("hera", 21)
		c.Incr("zeus", 6)
		c.Incr("ares", 1)
		c.Incr("athena", 4)
		wg.Done()
	}()
	go func() {
		c.Incr("zeus", 2)
		c.Incr("hera", 6)
		c.Incr("athena", 1)
		c.Incr("ares", 1)
		wg.Done()
	}()
	go func() {
		c.Incr("athena", 2)
		c.Incr("hera", 1)
		c.Incr("zeus", 3)
		c.Incr("ares", 1)
		wg.Done()
	}()
	// sync.WaitGroups: wait until all 4 threads report Done.
	// See the documentation.
	wg.Wait()

	expected := map[string]int{
		"zeus":   13,
		"hera":   29,
		"ares":   4,
		"athena": 8,
	}
	for k, v := range expected {
		if v != c.Get(k) {
			t.Errorf("counter %s: expected %d, got %d", k, v, c.Get(k))
		}
	}
}

func TestCounter(t *testing.T) {
	c := New()
	body(c, t)
}

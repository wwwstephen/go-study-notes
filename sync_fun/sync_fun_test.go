package syncfun

import (
	"sync"
	"testing"
	"time"
)

type Counter struct {
	Count int
}

func (c *Counter) Inc() {
	current := c.Count
	time.Sleep(time.Millisecond) // exaggerates the race condition
	c.Count = current + 1
}

func (c *Counter) Value() int {
	return c.Count
}

func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{}

		counter.Inc()
		counter.Inc()
		counter.Inc()

		if counter.Value() != 3 {
			t.Errorf("got %d, want %d", counter.Value(), 3)
		}
	})
}

func TestCounterConcurrency(t *testing.T) {
	counter := Counter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()

	if counter.Value() != 1000 {
		t.Errorf("got %d, want %d", counter.Value(), 1000)
	}
}

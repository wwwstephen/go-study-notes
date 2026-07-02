package mockingbird

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

const finalWord = "Go!"
const countdownStart = 3

const write = "write"
const sleep = "sleep"

// Sleeper defines something that can sleep (used for dependency injection)
type Sleeper interface {
	Sleep()
}

// SpyCountdownOperations acts as both a Sleeper and a Writer spy.
// It records the order of operations so we can assert behavior.
type SpyCountdownOperations struct {
	Calls []string
}

// Sleep records that Sleep was called
func (s *SpyCountdownOperations) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

// Write records that Write was called and simulates io.Writer
func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return len(p), nil
}

// Countdown prints 3..1 then "Go!", sleeping before each print
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(out, i)
	}
	fmt.Fprint(out, finalWord)
}

// real sleeper
// basically instead of writing to io.Writer, we sleep
type DefaultSleeper struct{}

func (d DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

// Test suite for Countdown
func TestCountdown(t *testing.T) {

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		// We only care about output here, so we use a real buffer
		Countdown(buffer, &SpyCountdownOperations{})

		got := buffer.String()
		want := "3\n2\n1\nGo!"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		spy := &SpyCountdownOperations{}

		// Use same spy for both interfaces (Writer + Sleeper behavior recorded)
		Countdown(spy, spy)

		want := []string{
			sleep, // before 3
			write,
			sleep, // before 2
			write,
			sleep, // before 1
			write,
			write, // final "Go!"
		}

		if !reflect.DeepEqual(want, spy.Calls) {
			t.Errorf("wanted calls %v got %v", want, spy.Calls)
		}
	})
}

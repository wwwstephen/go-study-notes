package dependencyinjection

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Chris")

	Greet(os.Stdout, "Peter")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

/*An example of dependency injection is writing output to an injected io.Writer rather than directly to os.Stdout.

By passing the dependency into your code, you can change its behavior depending on the environment. In production, you might pass os.Stdout so output appears in the terminal. In tests, you can pass a bytes.Buffer instead, allowing you to capture and verify the output without printing anything to the console.*/
package integers

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
func Test2(t *testing.T) {
	t.Run("my_test", func(t *testing.T) {
		if true != false {
			t.Errorf("Error in test")
		}
	})
}

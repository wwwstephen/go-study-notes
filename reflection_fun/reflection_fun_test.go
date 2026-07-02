package reflectionfun

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
}

func main() {
	x := Person{
		Name: "stephen",
	}

	fmt.Println(x.Name)
	v := reflect.ValueOf(x)
	// reflect is a reflection of an interface so
	// you can access the Fields directly.
	fmt.Println(v.Field(0))
}

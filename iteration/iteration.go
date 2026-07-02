package iteration

import "strings"

func Repeat(a string) string {
	var c strings.Builder; c.WriteString(a)
	for range 4 {
		c.WriteString(a)
	}
	
	return c.String()
}

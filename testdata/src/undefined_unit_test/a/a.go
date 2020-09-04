package a

import "fmt"

func f(n int) string { // want "This function's unit test is not defined"
	if n == 0 {
		return "zero"
	} else {
		return "not zero"
	}
}

func F2(s string) string { // OK
	return fmt.Sprintf("f2%s", s)
}
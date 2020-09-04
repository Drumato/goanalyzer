package in_if

import "fmt"

func f() {
	if x := f2(); x != 0 { // OK
		y := f2() // OK
		fmt.Println(x, y)
	} else {
		fmt.Println(0)
	}
}

func f2() int {
	return 0
}
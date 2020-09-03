package a

import "fmt"

func f() {
	x := f2() // want "This identifier is only referenced in a scope so should move the declaration to it"

	if true {
		y := f2() // OK
		fmt.Println(x, y)
	} else {
		fmt.Println(0)
	}
}

func f2() int {
	f()
	return 30
}

package in_for

import "fmt"

func f() {
	var i int // want "This identifier is only referenced in a scope so should move the declaration to it"
	for i = 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func f2() {
	var i int // OK
	for i = 0; i < 10; i++ {
		fmt.Println(i)
	}

	fmt.Println(i)
}
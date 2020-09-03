package closure

import "fmt"

func f() {
	x := 3 // OK

	func() {
		fmt.Println(x)
	}()

	fmt.Println(x)
}

func f2() {
	x := 3 // want "This identifier is only referenced in a scope so should move the declaration to it"

	func() {
		fmt.Println(x)
	}()

}
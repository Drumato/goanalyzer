package nested

import "fmt"

func f3() {
	outer := 1 // want "This identifier is only referenced in a scope so should move the declaration to it"

	if true {
		inner1 := 2 // want "This identifier is only referenced in a scope so should move the declaration to it"
		if true {
			inner2 := 3 // want "This identifier is only referenced in a scope so should move the declaration to it"
			if true {
				inner3 := 4 // OK

				fmt.Println(outer, inner1, inner2, inner3)
			} else {

			}
		} else {

		}
	} else {
		f4()
	}
}

func f4(){
	f3()
}
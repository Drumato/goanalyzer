package a

import "fmt"

func f(){
	x := f2() // want "NG"

	if true {
		y := f2() // OK
		fmt.Println(x, y)
	} else {
		fmt.Println(0)
	}
}

func f2() int {
	return 30
}
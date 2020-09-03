package samepkg

import "fmt"

func b() {
	fmt.Println(strGV1)
	fmt.Println(strGC1)

	var a strGT1 = "A"
	var b strGT2 = "A"
	fmt.Println(a, b)
}

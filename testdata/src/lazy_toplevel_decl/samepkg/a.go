package samepkg

import "fmt"

var (
	strGV1 string = "A" // OK
)

const (
	strGC1 string = "A" // OK
)

type strGT1 string   // OK
type strGT2 = string // OK

func a() {
	fmt.Println(strGV1)
	fmt.Println(strGC1)

	var a strGT1 = "A"
	var b strGT2 = "A"
	fmt.Println(a, b)
}

func f() {
	fmt.Println("B")
}

func f2() {
	a()
}

package a

import "fmt"

var (
	strGV1 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"
)

const (
	strGC1 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"
)

type strGT1 string   // want "This identifier is only referenced by a function so should move the declaration to it"
type strGT2 = string // want "This identifier is only referenced by a function so should move the declaration to it"

func A() {
	fmt.Println(strGV1)
	fmt.Println(strGC1)

	var a strGT1 = "A"
	var b strGT2 = "A"
	fmt.Println(a, b)
}

func B() {
	fmt.Println("B")
}

func C() {
	A()
}

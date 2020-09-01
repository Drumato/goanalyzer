package a

import "fmt"

var (
	strGV1 string = "A" // want "This identifier is declared in an unnecessarily wide scope"
)

const (
	strGC1 string = "A" // want "This identifier is declared in an unnecessarily wide scope"
)

type strGT1 string // want "This identifier is declared in an unnecessarily wide scope"
type strGT2 = string // want "This identifier is declared in an unnecessarily wide scope"

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
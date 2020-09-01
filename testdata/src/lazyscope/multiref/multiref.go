package multiref

import "fmt"

var (
	strGV1 string = "A" // want "This identifier is declared in an unnecessarily wide scope"
)
var strGV2 string = "A" // want "This identifier is declared in an unnecessarily wide scope"

const (
	strGC1 string = "A" // want "This identifier is declared in an unnecessarily wide scope"
)

const strGC2 string = "A" // want "This identifier is declared in an unnecessarily wide scope"

type strGT1 string// want "This identifier is declared in an unnecessarily wide scope"
type strGT2 = string // want "This identifier is declared in an unnecessarily wide scope"

type (
	strGT3 string// want "This identifier is declared in an unnecessarily wide scope"
	strGT4 = string // want "This identifier is declared in an unnecessarily wide scope"
)

func A() {
	fmt.Println(strGV1)
	fmt.Println(strGV2)
	fmt.Println(strGC1)
	fmt.Println(strGC2)

	var a strGT1 = "A"
	var b strGT2 = "A"
	var c strGT3 = "A"
	var d strGT4 = "A"

	fmt.Println(strGV1)
	fmt.Println(strGV2)
	fmt.Println(strGC1)
	fmt.Println(strGC2)

	var a2 strGT1 = "A"
	var b2 strGT2 = "A"
	var c2 strGT3 = "A"
	var d2 strGT4 = "A"

	fmt.Println(a, c, b+d, a2, c2, b2 + d2)
}
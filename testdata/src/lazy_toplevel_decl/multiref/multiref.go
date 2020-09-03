package multiref

import "fmt"

var (
	strGV1 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"
)
var strGV2 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"

const (
	strGC1 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"
)

const strGC2 string = "A" // want "This identifier is only referenced by a function so should move the declaration to it"

type strGT1 string   // want "This identifier is only referenced by a function so should move the declaration to it"
type strGT2 = string // want "This identifier is only referenced by a function so should move the declaration to it"

type (
	strGT3 string   // want "This identifier is only referenced by a function so should move the declaration to it"
	strGT4 = string // want "This identifier is only referenced by a function so should move the declaration to it"
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

	fmt.Println(a, c, b+d, a2, c2, b2+d2)
}

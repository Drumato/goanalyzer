package fp

import "fmt"

var (
	strGV1 string = "A" // OK
)
var strGV2 string = "A" // OK

const (
	strGC1 string = "A" //OK
)

const strGC2 string = "A" // OK

type strGT1 string   // OK
type strGT2 = string // OK

type (
	strGT3 string   // OK
	strGT4 = string // OK
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
	fmt.Println(a, c, b+d)
}

func B() {
	fmt.Println(strGV1)
	fmt.Println(strGV2)
	fmt.Println(strGC1)
	fmt.Println(strGC2)

	var a strGT1 = "A"
	var b strGT2 = "A"
	var c strGT3 = "A"
	var d strGT4 = "A"
	fmt.Println(a, c, b+d)
}

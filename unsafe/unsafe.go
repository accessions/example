package main

import (
	"fmt"
	"unsafe"
)

func main()  {
	var i int32 = 1111111111
	var ii int64 = 1111111111111111111
	bools := unsafe.Sizeof(true)
	ints := unsafe.Sizeof(1111111)
	fmt.Println("int", ints)
	fmt.Println("bool ",bools)
	int3 := unsafe.Sizeof(i)
	fmt.Println("int32", int3)
	int6 := unsafe.Sizeof(ii)
	fmt.Println("int64", int6)
	var s = "ab"
	st := unsafe.Sizeof(s)
	fmt.Println("string", st)
	var f32 float32 = 1.42
	var f64 float64 =  1.43
	f := unsafe.Sizeof(1.42)
	fmt.Println("float", f)
	f3 := unsafe.Sizeof(f32)
	fmt.Println("float32", f3)
	f6 := unsafe.Sizeof(f64)
	fmt.Println("float64", f6)
	v := struct {
		Name string
		Age int32
	}{Name: "11111",Age: 30}
	str := unsafe.Sizeof(v)
	fmt.Println("struct", str)

}

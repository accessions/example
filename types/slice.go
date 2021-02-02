package main

import (
	"fmt"
	"github.com/accessions/core/utils"
)


func main()  {
	a := []interface{}{"A","B","C"}
	x, a := a[0], a[1:]
	fmt.Println(x,a)
	shift, i := utils.Shift(a)
	fmt.Println(shift, i)



}



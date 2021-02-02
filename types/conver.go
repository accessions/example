package main

import "fmt"

func conversion ()  {
	iface := []interface{}{1,"2","3"}
	fmt.Println(iface...)
	strs := []string{"1","A","B"}
	//fmt.Println(strs...)
	fmt.Println(strs)
}



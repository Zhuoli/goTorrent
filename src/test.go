package main

import (
	"fmt"
	"sync"
)

type ST struct{
	val	string
}


func main() {
	
	var tmp sync.Pool
	el:=tmp.Get()
	fmt.Println(el)
}
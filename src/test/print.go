package test

import (
	"strings"
	"fmt"
)

func PrintMap(m map[string][]string){
		for k,v:=range m{
		fmt.Println("Key: "+k+" values: "+strings.Join(v,","))
	}
}
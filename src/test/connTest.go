package main

import (
	"fmt"
	"conn"
)

func httpGetTest(){
	var urlStr string="http://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol"
	http :=conn.GetHttp(urlStr)
	http.Get("/wiki/Hypertext_Transfer_Protocol",0,1)
	
}

func main(){
	fmt.Println("Hello world")
	httpGetTest()
}
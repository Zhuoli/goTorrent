package test

import (
	"fmt"
	"conn"
	"strconv"
	"strings"
)

func httpGetTest(){
	var urlStr string="http://download.nextag.com/apache/maven/maven-3/3.2.5/binaries/apache-maven-3.2.5-bin.tar.gz"
	strs:=strings.Split(urlStr,"/")
	http :=conn.GetHttp(urlStr)
	http.Get("/wiki/Hypertext_Transfer_Protocol",0,1)
	response:=http.Response()
	fmt.Println("length: "+strconv.Itoa(http.GetContentLength(response)))
	fmt.Println(fmt.Sprintf("Is accept range: %t",http.IsAcceptRange(response)))
	con:=&conn.CONN{}
	con.Get(http,urlStr,strs[len(strs)-1])
	
}

func main(){
	httpGetTest()
	
}
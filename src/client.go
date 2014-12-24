package main

import (
	"fmt"
	"strings"
//	"thread"
	"conn"
)

const(
	DEBUG bool=true
)

func downloadFromUrl(url string) {
	fileName := getFileName(url)
	fmt.Println("Downloading", url, "to", fileName)
	conn:=conn.GetConn(url)
	if  conn.IsAcceptRange{
		c:=make(chan int)
		go conn.WriteToFile(fileName,c)
		<-c
	}else{
		fmt.Println("target url doesn't accept range")
		c:=make(chan int)
		go conn.WriteToFile(fileName,c)
		<-c
	}
	fmt.Println("DONE")
}

func getFileName(url string) (name string){
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

func main() {
//	url := "http://shakespeare.mit.edu/lll/full.html"
	url :="http://upload.wikimedia.org/wikipedia/commons/2/2f/Space_Needle002.jpg"
//	url := "http://mirrors.sonic.net/apache/maven/maven-3/3.2.5/binaries/apache-maven-3.2.5-bin.tar.gz"
//	url :="http://www.ccs.neu.edu/course/cs5500f14/policies.html"
	downloadFromUrl(url)
}
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
	conn:=conn.InitConn(url)
	if  conn.IsAcceptRange && conn.HasContentLength{
		c:=make(chan int)
		go conn.WriteToFile(fileName,c)
//		thread.RangeDownload(fileName,conn,c)
		<-c
	}else{
		fmt.Println("target url doesn't accept range")
		c:=make(chan int)
		go conn.WriteToFile(fileName,c)
		fmt.Println(<-c)
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
//	url :="http://upload.wikimedia.org/wikipedia/commons/2/2f/Space_Needle002.jpg"
	url := "http://mirrors.sonic.net/apache/maven/maven-3/3.2.5/binaries/apache-maven-3.2.5-bin.tar.gz"
//	url :="http://www.ccs.neu.edu/course/cs5500f14/policies.html"
	downloadFromUrl(url)
}
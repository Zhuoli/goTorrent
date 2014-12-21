package main

import (
	"fmt"
	"strings"
	"thread"
)

const(
	DEBUG bool=true
)

func downloadFromUrl(url string) {
	fileName := getFileName(url)
	fmt.Println("Downloading", url, "to", fileName)
	conn:=thread.GetConn(url)
	if(DEBUG){
		fmt.Println(fmt.Sprintf("Content length: %d bytes",conn.GetLength()))
		fmt.Println(fmt.Sprintf("Accept range: %t",conn.GetIsAcceptRange()))
	}
	if false && conn.GetIsAcceptRange(){
		
	}else{
		fmt.Println("target url doesn't accept range")
		c:=make(chan int)
		go thread.Read2File(conn,c)
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
	downloadFromUrl(url)
}
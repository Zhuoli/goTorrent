package main

import (
	"fmt"
	"net/http"
	"strings"
	"thread"
	"strconv"
)

func downloadFromUrl(url string, threadsNum int) {
	fileName := getFileName(url)
	fmt.Println("Downloading", url, "to", fileName)
	fmt.Println("threadsNum: "+strconv.Itoa(threadsNum))
	response,err:=http.Get(url)
	if err!=nil{
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	length:=response.ContentLength
	fmt.Println("contentlength : " + strconv.Itoa(int(length))+"byte")
	block:=int64(0)
	if length%int64(threadsNum)!=0{
		block=length/int64(threadsNum)+int64(threadsNum)
	}else{
		block=length/int64(threadsNum)
	}
	fmt.Println("block size: " + strconv.FormatInt(block,10)+"byte")
	c:=make(chan string,threadsNum)
	for i:=0;i<threadsNum;i++{
		go thread.DownLoadThread(c,url,fileName,int64(i),block)
	}
	
	for i:=0;i<threadsNum;i++{
		fmt.Println(<-c)
	}
	 
	fmt.Println("content length: "+strconv.Itoa(int(response.ContentLength/1024))+"KB")
}

func getFileName(url string) (name string){
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

func main() {
	url := "http://download.nextag.com/apache/maven/maven-3/3.2.3/binaries/apache-maven-3.2.3-bin.tar.gz"
	downloadFromUrl(url,5)
}
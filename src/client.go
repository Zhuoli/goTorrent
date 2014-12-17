package main

import (
	"fmt"
	"net/http"
	"strings"
	"thread"
	"strconv"
	"os"
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
//	response.Header.Get
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
	// open the file, if not exist create it
	File,err:=os.Create(fileName)
	buf:=make([]byte, length)
	if err!=nil{	
		panic(err)
		}
	defer File.Close()
	for i:=0;i<threadsNum;i++{
		go thread.DownLoadThread(c,buf,url,fileName,int64(i),block)
	}
	
	for i:=0;i<threadsNum;i++{
		fmt.Println(<-c)
	}
	_,err=File.Write(buf)
	if err!=nil{
		panic(err)
	}
	fmt.Println("content length: "+strconv.Itoa(int(response.ContentLength/1024))+"KB")
}

func getFileName(url string) (name string){
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

func main() {
	url := "http://upload.wikimedia.org/wikipedia/commons/2/2f/Space_Needle002.jpg"
	downloadFromUrl(url,5)
}
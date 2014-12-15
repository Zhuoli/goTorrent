package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadFromUrl(url string) {
	fileName := getFileName(url)
	fmt.Println("Downloading", url, "to", fileName)
	
	client :=&http.Client{}
	req,_:=http.NewRequest("GET",url,nil)
	req.Header.Set("Range","bytes=start-end")
}

func getFileName(url string) name String{
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

func main() {
	url := "http://download.nextag.com/apache/maven/maven-3/3.2.3/binaries/apache-maven-3.2.3-bin.tar.gz"
	downloadFromUrl(url)
}
package main

import (
	"fmt"
	"strings"
	"flag"
	"os"
	"conn"
	"bitTorrent/bencode"
)

const(
	DEBUG bool=false
)
// http download
func downloadFromUrl(url string) {
	fileName := getFileName(url)
	fmt.Println("Downloading", url, "to", fileName)
	conn:=conn.InitConn(url)
	c:=make(chan int)
	go conn.WriteToFile(fileName,c)
	fmt.Println(<-c)
	fmt.Println("DONE")
}

//BitTorrent seed download
func downloadBitTorrent(seed string){
	fmt.Println("torrent file: "+seed)
	f,err:=os.Open(seed)
	if err!=nil{
		panic(err)
	}
	defer f.Close()
	bencode.Decode(f)
}

func getFileName(url string) (name string){
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

func main() {
//	url := "http://shakespeare.mit.edu/lll/full.html"
//	url :="http://upload.wikimedia.org/wikipedia/commons/2/2f/Space_Needle002.jpg"
//	url := "http://mirrors.sonic.net/apache/maven/maven-3/3.2.5/binaries/apache-maven-3.2.5-bin.tar.gz"
//	url :="http://www.ccs.neu.edu/course/cs5500f14/policies.html"
	url:=flag.String("url","","url address")
	torrent:=flag.String("torrent","","bit-torrent seed")
	flag.Parse()
	if len(*url)!=0 {
		downloadFromUrl(*url)
	}else if len(*torrent)!=0{
		downloadBitTorrent(*torrent)
	}else{
		fmt.Println("Error usage: ./client -url=")
		fmt.Println("Error usage: ./client -torrent=")
		return
	}
}
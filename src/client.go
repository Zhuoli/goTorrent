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
	
	// DEBUG idenciate
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

/*
	Get file name by parsing the given Url
	Given: Url
	Returns: File name
*/
func getFileName(url string) (name string){
	tokens :=strings.Split(url,"/")
	fileName:=tokens[len(tokens)-1]
	return fileName
}

/**
	Client main method.
	* Currently suport http downloading.
	Usage: 
	./clien -url=XX
	./client -torrent=xx
	Downloaded files will be stored under ./download
**/
func main() {
	
	// Sets url parameter format 
	url:=flag.String("url","","url address")
	
	// Sets torrent parameter format
	torrent:=flag.String("torrent","","bit-torrent seed")
	
	// Parses the argument inputs from command line
	flag.Parse()
	
	// Validate inputs
	if len(*url)!=0 {
		
		// Download from http url server
		downloadFromUrl(*url)
	}else if len(*torrent)!=0{
		
		// Download from BitTorrent seed
		downloadBitTorrent(*torrent)
	}else{
		
		// Error usage
		fmt.Println("Error usage: ./client -url=")
		fmt.Println("Error usage: ./client -torrent=")
		return
	}
}
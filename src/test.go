package main

import (
	"fmt"
	"os"
	"strconv"
)

func test(c chan string, i int,fileName string) {
	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	bytes:=[]byte{'h','e','l','l','o',' ','w','o','r','l','d'}

	n,err:=output.Write(bytes)
	if err!=nil{
		panic(err)
	}
	c<-"Thread "+strconv.Itoa(i)+" "+strconv.Itoa(n)+ " bytes wrote."
}

func main() {
	threadsNum:=5
	c:=make(chan string,threadsNum)
	for i:=0;i<threadsNum;i++{
		go test(c,i,"hello.txt")
	}
	for i:=0;i<threadsNum;i++{
		fmt.Println(<-c)
	}
}
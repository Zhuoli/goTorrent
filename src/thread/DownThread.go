package thread

import (
	"net/http"
	"fmt"
	"strconv"
	"os"
	"io"
)

func DownLoadThread(c chan string,url,filename string, i,block int){
	start:=0
	end:=0
	//Set a header to a request first. Pass the request to a client.
	client := &http.Client{}
	req, err := http.NewRequest("GET", url,nil)
	if err!=nil{
		fmt.Println(err)
	}
	start=i*block
	end=(i+1)*block
	req.Header.Set("Range","bytes="+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"")
	resp, err := client.Do(req)
	if err!=nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	
	// TODO: check file existence first with io.IsExist
	// return a File object
	output, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error while creating", filename, "-", err)
		return
	}
	defer output.Close()
	
	n, err := io.Copy(output, resp.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	s:="Thread "+strconv.Itoa(i)+" Done, received: "+strconv.Itoa(int(n))+" bytes"
	c <- s
}
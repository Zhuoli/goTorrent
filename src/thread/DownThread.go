package thread

import (
	"net/http"
	"fmt"
	"strconv"
	"os"
)

func DownLoadThread(c chan string,url,filename string, i,block int64){

	//Set a header to a request first. Pass the request to a client.
	client := &http.Client{}
	req, err := http.NewRequest("GET", url,nil)
	if err!=nil{
		fmt.Println(err)
	}
	start:=i*block
	end:=(i+1)*block
	// set http header filed: key = "Range" value = "..."
	req.Header.Set("Range","bytes="+strconv.FormatInt(start,10)+"-"+strconv.FormatInt(end,10)+"")
	resp, err := client.Do(req)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	
	// open the file, if not exist create it
	output,err:=os.Open(filename)
	if err!=nil{
		output,err=os.Create(filename)
		if err!=nil{
			fmt.Println(err)
			return
		}
	}
	defer output.Close()
	
//	n, err := io.Copy(output, resp.Body)
	var bytes = make([] byte,end-start)
	n,err:=resp.Body.Read(bytes)
	if err != nil {
		fmt.Println("Error while reading bytes", err)
		return
	}
	if output==nil{
		fmt.Println("output is nil")
	}
	n,err=output.WriteAt(bytes,start)
	if err != nil {
		fmt.Println("Error while writting", url, "-", err)
		return
	}
	s:="Thread "+strconv.FormatInt(i,10)+" Done, received: "+strconv.Itoa(n)+" bytes"
	c <- s
}
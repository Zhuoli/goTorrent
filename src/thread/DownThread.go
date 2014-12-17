package thread

import (
	"net/http"
	"fmt"
	"strconv"
)

func DownLoadThread(c chan string,buf []byte,url,filename string, i,block int64){

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
	

	var bytes = make([] byte,end-start)
	n,err:=resp.Body.Read(bytes)
	if err != nil {
		fmt.Println("Error while reading bytes", err)
		return
	}
	copy(buf[start:end],bytes[:])
	s:="Thread "+strconv.FormatInt(i,10)+" Done, received: "+strconv.Itoa(n)+" bytes"
	c <- s
}


package conn

import (
	"fmt"
	"strings"
	"strconv"
	"os"
)
type CONN struct{
	http	*HTTP
	IsAcceptRange	bool
	HasContentLength bool
	Content_length	int64
	IsTransfer_Encoding bool
}
const (
    DEBUG=true
)

func GetConn(url string)*CONN{
	conn:=&CONN{}
	http:=GetHttp(url)
	defer http.conn.Close()
	conn.http=http
	if http.connect()==false{return nil}
	http.Get(url,0,0)
	response:=http.Response()
	if strings.Contains(response,"Accept-Ranges:"){
		conn.IsAcceptRange=true
	}
	if strings.Contains(response,"Content-Length:"){
		conn.HasContentLength=true
		index:=strings.Index(response,"Content-Length:")
		substr:=response[index+len("Content-Length: "):]
		index=strings.Index(substr,"\n")
		intstr:=substr[:index]
		length,err:=strconv.ParseInt(intstr,10,64)
		conn.Content_length=length
		if err!=nil{
			panic("Error, incorrect content-length"+intstr)
		}
	}else if strings.Contains(response,"Transfer-Encoding: chunked"){
		conn.IsTransfer_Encoding=true
	}
	if DEBUG{
		fmt.Println(fmt.Sprintf("Accept-Range: %t",conn.IsAcceptRange))
		if conn.HasContentLength{
			fmt.Println(fmt.Sprintf("Content-length: %d",conn.Content_length))
		}else if conn.IsTransfer_Encoding{
			fmt.Println("Transfer_Encoding: bytes")
		}else{
			fmt.Println(response)
		}
	}
	return conn
	
}

func (this *CONN)WriteToFile(fileName string, c chan int){
	src,err:=os.Stat("./download")
	if err!=nil || !src.IsDir(){
		os.Mkdir("./download",0777)
	}
	if this.HasContentLength{
		this.http.WriteToFileContentLength(fileName,this.Content_length)
	}else if this.IsTransfer_Encoding{
		this.http.WriteToFileTruncked(fileName)
	}
	c<-1
}
package conn

import (
	"fmt"
	"strings"
	"strconv"
)
type CONN struct{
	http	*HTTP
	IsAcceptRange	bool
	HasContentLength bool
	Content_length	int64
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
	}
	if DEBUG{
		fmt.Println(fmt.Sprintf("Accept-Range: %t",conn.IsAcceptRange))
		if conn.HasContentLength{
			fmt.Println(fmt.Sprintf("Content-length: %d",conn.Content_length))
			}else{
			fmt.Println("Don't has Content-length")
		}
	}
	return conn
	
}

func (this *CONN)WriteToFile(fileName string, c chan int){
	this.http.WriteToFile(fileName,this.Content_length)
	c<-1
}
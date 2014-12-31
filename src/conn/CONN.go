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
    DEBUG=false
    threadNum=4
)

func InitConn(url string)*CONN{
	conn:=&CONN{}
	http:=GetHttp(url)
	conn.http=http
	netconn,err:=http.connect()
	defer netconn.Close()
	if err!=nil{return nil}
	http.Get(netconn,0,-1)
	response:=http.Response(netconn)
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
		this.RangeDownload(fileName)
	}else if this.IsTransfer_Encoding{
		this.http.WriteToFileTruncked(fileName)
	}
	c<-1
}

func (this *CONN)RangeDownload(filename string){
	blocksize:=this.Content_length/threadNum+threadNum
	receiver:=make([]chan int,threadNum)
	for i:=0;i<threadNum;i++{
		receiver[i]=make(chan int)
		http:=&HTTP{
			Protocol:this.http.Protocol,
			Host:this.http.Host,
			Path:this.http.Path,
			Port:this.http.Port,
			UserAgent:this.http.UserAgent,
		}
		end:=int64(i+1)*blocksize
		if end>this.Content_length{
			end=this.Content_length
		}
		go http.WriteToFileContentLength(receiver[i],filename,int64(i)*blocksize,end)
	}
	
	for i:=0;i<threadNum;i++{
		<-receiver[i]
	}
}

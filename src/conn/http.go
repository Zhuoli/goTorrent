package conn

import (
	"fmt"
	"net"
	"net/url"
	"io"
	"os"
    "regexp"
    "strings"
    "strconv"
)
const (
	BUFFER_SIZE=1024
)
type HTTP struct{
	Protocol	string
	Host 		string
	Path		string
	Port		int
	Header		string
	UserAgent	string
}

func GetHttp(strUrl string)(*HTTP){
	http :=&HTTP{}
	u,err:=url.Parse(strUrl)
	if err!=nil{
		fmt.Println("ERROR:", err.Error)
		return nil
	}
	http.Protocol=u.Scheme
	http.Host=u.Host
	http.Path=u.Path
	http.Port=80
	http.Header=""
	http.UserAgent="Mozilla/5.0"
	return http
}

func (this *HTTP)AddHeader(header string){
	this.Header+=header+"\r\n"
}

func (this *HTTP) connect()(conn net.Conn, e error){
    address := fmt.Sprintf("%s:%d", this.Host, this.Port)
    conn, Error := net.Dial("tcp", address)
    if Error != nil {
        fmt.Println("ERROR: ", Error.Error())
        return nil,Error
    }
    return conn,nil
}
/* Get method for http */
func (this *HTTP)Get(conn net.Conn,range_from, range_to int64){
	this.AddHeader(fmt.Sprintf("GET %s HTTP/1.1",this.Path))
	this.AddHeader(fmt.Sprintf("Host: %s ",this.Host))
	if range_to==-1{
		this.AddHeader(fmt.Sprintf("Range: bytes=%d-",range_from))
	}else{
		this.AddHeader(fmt.Sprintf("Range: bytes=%d-%d",range_from,range_to))
	}
	this.AddHeader(fmt.Sprintf("User-Agent: %s",this.UserAgent))
	this.AddHeader("Connection: keep-alive")
	this.AddHeader("\r\n")
	_,Error:=conn.Write([]byte(this.Header))
	if Error!=nil{
		fmt.Println("Error: ",Error.Error())
	}
}

func (this *HTTP) Response(conn net.Conn)string{
	var headerResponse string
//	defer this.conn.Close()
	data :=make([]byte,1)
	for i:=0;;{
		n,err:=conn.Read(data)
		if err!=nil{
			if err!=io.EOF{
				Error:=err
				fmt.Println("ERROR: ",Error.Error)
				return "";
			}
		}
		if data[0]=='\r'{
			continue
		}else if data[0]=='\n'{
			if i==0{
				break
			}
			i=0
		}else{
			i++
		}
		headerResponse+=string(data[:n])
	}
//	if !strings.Contains(headerResponse,"HTT"){
//		fmt.Println("ERROR: header response not OK")
//		fmt.Println(headerResponse)
//		panic(errors.New(headerResponse))
//	}
	return headerResponse
}

func (this *HTTP) GetContentLength(headerResponse string) int {
    ret := 0
    r, err := regexp.Compile(`Content-Length: (.*)`)
    if err != nil {
        fmt.Println("ERROR: ", err.Error())
        return ret
    }
    result := r.FindStringSubmatch(headerResponse)
    if len(result) != 0 {
        s := strings.TrimSuffix(result[1], "\r")
        ret, _ = strconv.Atoi(s)
    }
    return ret
}
func (this *HTTP) IsAcceptRange(headerResponse string) bool {
    ret := false

    if strings.Contains(headerResponse, "Content-Range") || 
        strings.Contains(headerResponse, "Accept-Ranges"){
        ret = true
    }

    return ret
}
func (this *HTTP) WriteToFileContentLength(c chan int,outputFileName string,start,end int64) {
//	fmt.Println("WriteToFileContentLength")
	conn,_:=this.connect()
	defer conn.Close()
	this.Get(conn,start,end)
	_=this.Response(conn)
    f, err := os.OpenFile("./download/"+outputFileName, os.O_CREATE | os.O_WRONLY, 0664)
    defer f.Close()
    if err != nil { panic(err) }
    data := make([]byte, BUFFER_SIZE)
    offset:=start
    for{
    	n,err:=conn.Read(data)
    	if err!=nil{
    		if err!=io.EOF{
    			Error:=err
    			panic(Error.Error())
    			return
    		}
    	}
    	f.WriteAt(data[:n],int64(offset))
    	offset+=int64(n)
    	if err==io.EOF{break}	
    	if int64(offset)>=end{break}
    }
    fmt.Println(fmt.Sprintf("Write from %d to %d, Received %d bytes",start,end,offset-start))
	c<-1
    return
}
func (this *HTTP) WriteToFileTruncked(fileName string){
	conn,_:=this.connect()
	this.Get(conn,0,-1)
	_=this.Response(conn)
//	fmt.Println(res)
    f, err := os.OpenFile("./download/"+fileName, os.O_CREATE | os.O_WRONLY, 0664)
    defer f.Close()
    if err != nil { panic(err) }
    fileoffset:=0
    // chunks
    for {
    	size:=this.getChunckSize(conn)
    	if size==0{
    		break
    	}
    	// offset in chunk
    	chunkoffset:=0
    	for{
	   	    data := make([]byte, 1)
	    	n,err:=conn.Read(data)
	    	if err!=nil{
	    		if err!=io.EOF{
	    			Error:=err
	    			panic(Error.Error())
	    			return
	    		}
	    	}
	    	//reach the end of this chunk
	    	if chunkoffset+n>=size {
	    		f.WriteAt(data[:n],int64(fileoffset+chunkoffset))
		    	conn.Read(data)
		    	conn.Read(data)
	    		break
	    	}else{
	    		f.WriteAt(data[:n],int64(fileoffset+chunkoffset))
	    	}
	    	if err==io.EOF{return}
	    	chunkoffset+=n
    	}
	    fileoffset+=size
    }
    fmt.Println(fmt.Sprintf("Received %d bytes",fileoffset))
    return
}
func (this *HTTP)getChunckSize(conn net.Conn) int{
	expr:=""
	for{
		data:=make([]byte,1)
		n,err:=conn.Read(data)
		if err != nil {
            if err != io.EOF {
                Error := err
                fmt.Println("ERROR:", Error.Error())
                return 0
            }
        }
		if data[0] == '\r' {
            continue
        } else if data[0] == '\n' {
        	break
        }
		expr+=string(data[:n])
		
	}
	num,_:=strconv.ParseInt(expr,16,32)
	return int(num)
}
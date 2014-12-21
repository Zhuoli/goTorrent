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

type HTTP struct{
	Protocol	string
	Host 		string
	Port		int
	Header		string
	UserAgent	string
	Error		error
	conn 		net.Conn
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
	http.Port=80
	http.Header=""
	http.UserAgent="GoTOrrent"
	http.connect()
	return http
}

func (this *HTTP)AddHeader(header string){
	this.Header+=header+"\r\n"
}

func (this *HTTP) connect()bool{
    address := fmt.Sprintf("%s:%d", this.Host, this.Port)
    this.conn, this.Error = net.Dial("tcp", address)
    if this.Error != nil {
        fmt.Println("ERROR: ", this.Error.Error())
        return false
    }
    return true
}
/* Get method for http */
func (this *HTTP)Get(url string, range_from, range_to int){
	this.AddHeader(fmt.Sprintf("GET %s HTTP/1.1",url))
	this.AddHeader("Connection: close")
	this.AddHeader(fmt.Sprintf("Host: %s ",this.Host))
	this.AddHeader(fmt.Sprintf("Range: bytes=%d-%d",range_from,range_to))
	this.AddHeader(fmt.Sprintf("User-Agent: %s",this.UserAgent))
	this.AddHeader("")
	_,this.Error=this.conn.Write([]byte(this.Header))
	if this.Error!=nil{
		fmt.Println("Error: ",this.Error.Error())
	}
}

func (this *HTTP) Response()string{
	var headerResponse string
	defer this.conn.Close()
	data :=make([]byte,1)
	for i:=0;;{
		n,err:=this.conn.Read(data)
		if err!=nil{
			if err!=io.EOF{
				this.Error=err
				fmt.Println("ERROR: ",this.Error.Error)
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
	return headerResponse
}

func (this *HTTP) GetContentLength(headerResponse string) int {
    ret := 0
    r, err := regexp.Compile(`Content-Length: (.*)`)
    if err != nil {
        this.Error = err
        fmt.Println("ERROR: ", err.Error())
        return ret
    }
    result := r.FindStringSubmatch(headerResponse)
    if len(result) != 0 {
        s := strings.TrimSuffix(result[1], "\r")
        ret, this.Error = strconv.Atoi(s)
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
func (this *HTTP) WriteToFile(outputFileName string, old_range_from int, chunkSize,buffer_size int) {
	offset:=chunkSize
    defer this.conn.Close()
    resp := ""
    data := make([]byte, 1)
    for i := 0; ; {
        data := make([]byte, 1)
        n, err := this.conn.Read(data)
        if err != nil {
            if err != io.EOF {
                this.Error = err
                fmt.Println("ERROR:", this.Error.Error())
                return
            }
        }
        if data[0] == '\r' {
            continue
        } else if data[0] == '\n' {
            if i == 0 {
                break
            }
            i = 0
        } else {
            i++
        }
        resp += string(data[:n])
    }
//    chunkName := fmt.Sprintf("%s.part.%d", outputFileName, old_range_from)
    f, err := os.OpenFile(outputFileName, os.O_CREATE | os.O_WRONLY, 0664)
    defer f.Close()
    if err != nil { panic(err) }
    data = make([]byte, buffer_size)
    for {
        n, err := this.conn.Read(data)
        if err != nil {
            if err != io.EOF {
                this.Error = err
                fmt.Println("ERROR:", this.Error.Error())
                return
            }
        }
        f.WriteAt(data[:n], int64(chunkSize))
        offset += n
        if err == io.EOF { return }
    }
    return
}
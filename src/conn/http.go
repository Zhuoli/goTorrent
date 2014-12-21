package conn

import (
	"fmt"
	"net"
	"net/url"
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
func (this *HTTP)Get(url string, range_from, range_to int64){
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

func (this *HTTP) Response(){
	
}
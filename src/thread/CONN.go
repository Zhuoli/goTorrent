package thread

import (
	"fmt"
	"net/http"
	"test"
)
type CONN struct{
	length		int64
	isAcceptRange	bool
	url			string
}

func GetConn(url string) *CONN{
	conn :=&CONN{url:url}
	response,err:=http.Get(url)
	if err!=nil{
		fmt.Println("Error while downloading", url, "-", err)
		return nil
	}
	defer response.Body.Close()
	if len(response.Header.Get("Accept-Ranges"))!=0{
		conn.isAcceptRange=true
	}
	conn.length=response.ContentLength
	if err!=nil{
		panic(err)
	}
	test.PrintMap(response.Header)
//	fmt.Println("header.get Content-Length: "+response.Header.Get("Content-Length"))
	return conn
}

func(this *CONN)GetLength()int64{
	return this.length
}
func(this *CONN)GetIsAcceptRange() bool{
	return this.isAcceptRange
}


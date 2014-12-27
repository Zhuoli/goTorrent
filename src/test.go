package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func getHtmlHeader(url string)string{
	response,err:=http.Get(url)
	if err!=nil{
		fmt.Println("Error while downloading", url, "-", err)
		return "";
	}
	defer response.Body.Close()
	
	body, err := ioutil.ReadAll(response.Body)
	if err!=nil{
		panic(err)
	}
	return string(body);
	
}


func main() {
}
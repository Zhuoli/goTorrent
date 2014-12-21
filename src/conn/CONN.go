package conn

import (

)
type CONN struct{
	
}
const (
    buffer_size int = 102400
)

func (this *CONN) Get(http *HTTP,url,fileName string){
	if http.connect()==false{return}
	http.Get(url,0,0)
    http.WriteToFile(fileName, 0, 0,buffer_size)
	
	
}
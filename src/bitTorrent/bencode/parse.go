package bencode

import (
	"io"
	"sync"
	"bufio"
	"errors"
	"strconv"
)

type builder interface{
	Int64(i int64)
	Uint64(i uint64)
	String(s string)
	Array()
	Map()
	
	Elem(i int)builder
	Key(s string)builder
	
	Flush()
}

func collectInt(r *bufio.Reader,delim byte)(buf []byte,err error){
	for{
		var c byte
		c,err = r.ReadByte()
		if err!=nil{
			return
		}
		if c==delim{
			return
		}
		if !(c=='-' || (c>='0' && c<='9')){
			err =errors.New("expected digit")
			return
		}
		buf=append(buf,c)
	}
	return
}

func decodeInt64(r *bufio.Reader,delim byte)(data int64,err error){
	buf,err:=collectInt(r,delim)
	if err!=nil{
		return 
	}
	data,err=strconv.ParseInt(string(buf),10,64)
	return
}

func decodeString(r *bufio.Reader)(data string,err error){
	length,err:=decodeInt64(r,':')
	if err!=nil{
		return
	}
	if length<0{
		err=errors.New("Bad string length")
		return
	}
	var buf = make([]byte,length)
	_,err=io.ReadFull(r,buf)
	if err!=nil{
		return
	}
	data = string(buf)
	return
}
func parseFromReader(r *bufio.Reader, build builder) (err error) {
	c,err:=r.ReadByte()
	if err!=nil{
		goto exit
	}
	switch {
		case c>='0' && c<='9':
		err=r.UnreadByte()
		if err!=nil{
			goto exit
		}
		var str string
		str,err=decodeString(r)
		if err!=nil{
			goto exit
		}
		build.String(str)
	}
	
exit:
	build.Flush()
	return	
}
func parse(r io.Reader,builder builder)(err error){
	buf:=newBufioReader(r)
	// put the bufio.reader to pool
	defer bufioReaderPool.Put(buf)
	return parseFromReader(buf,builder)
}

// a Pool variable
var bufioReaderPool sync.Pool

func newBufioReader(r io.Reader) *bufio.Reader {
	//Get selects an arbitrary item from the Pool, removes it from the Pool, and returns it to the caller.
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}
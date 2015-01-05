package bencode

import (
	"io"
)

func Decode(r io.Reader)(data interface{},err error){
	jb :=newDecoder(nil,nil)
	err=parse(r,jb)
	if err==nil{
		data=jb.Copy()
	}
	return
}

type decoder struct{
	value interface{}
	container interface{}
	index interface{}
}

func newDecoder(container interface{},key interface{}) *decoder{
	return &decoder{container:container,index:key}
}


func (j *decoder) Int64(i int64) { j.value = int64(i) }

func (j *decoder) Uint64(i uint64) { j.value = uint64(i) }

func (j *decoder) Float64(f float64) { j.value = float64(f) }

func (j *decoder) String(s string) { j.value = s }

func (j *decoder) Bool(b bool) { j.value = b }

func (j *decoder) Null() { j.value = nil }

func (j *decoder) Array() { j.value = make([]interface{}, 0, 8) }

func (j *decoder) Map() { j.value = make(map[string]interface{}) }

func (j *decoder) Elem(i int) builder {
	v, ok := j.value.([]interface{})
	if !ok {
		v = make([]interface{}, 0, 8)
		j.value = v
	}
	/* XXX There is a bug in here somewhere, but append() works fine.
	lens := len(v)
	if cap(v) <= lens {
		news := make([]interface{}, 0, lens*2)
		copy(news, j.value.([]interface{}))
		v = news
	}
	v = v[0 : lens+1]
	*/
	v = append(v, nil)
	j.value = v
	return newDecoder(v, i)
}

func (j *decoder) Key(s string) builder {
	m, ok := j.value.(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		j.value = m
	}
	return newDecoder(m, s)
}

func (j *decoder) Flush() {
	switch c := j.container.(type) {
	case []interface{}:
		index := j.index.(int)
		c[index] = j.Copy()
	case map[string]interface{}:
		index := j.index.(string)
		c[index] = j.Copy()
	}
}

// Get the value built by this builder.
func (j *decoder) Copy() interface{} {
	return j.value
}
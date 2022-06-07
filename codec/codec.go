package codec

import "io"

type Header struct {
	ServiceMethod string //format "Service.Method"
	Seq           uint64 //请求序号，可以认为某个请求的ID，区分不同的请求
	Error         string
}

//对消息体进行编解码的接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

//下边的与工厂模式类似，返回的是构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
	//other types
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	//other types
}

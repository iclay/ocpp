package server

import (
	"unsafe"
)

type convert interface {
	StringToBytes(string) []byte
	BytesToString([]byte) string
}

type rawConvert struct{}

func (rawConvert) StringToBytes(s string) []byte {
	return []byte(s)
}
func (rawConvert) BytesToString(b []byte) string {
	return string(b)
}

//custom conversion can optimize memory and reduce unnecessary copies
type customConvert struct{}

func (customConvert) StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (customConvert) BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))

}

func SupportCustomConversion(support bool) opt {
	return func(o *option) {
		switch support {
		case true:
			o.convert = customConvert{}
		default:
			o.convert = rawConvert{}
		}
	}
}

func Bytes(s string) []byte {
	return options.convert.StringToBytes(s)
}

func String(b []byte) string {
	return options.convert.BytesToString(b)
}

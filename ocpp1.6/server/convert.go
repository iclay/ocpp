package server

import (
	"unsafe"
)

type convert interface {
	stringToBytes(string) []byte
	bytesToString([]byte) string
}

type rawConvert struct{}

func (*rawConvert) stringToBytes(s string) []byte {
	return []byte(s)
}
func (*rawConvert) bytesToString(b []byte) string {
	return string(b)
}

//custom conversion can optimize memory and reduce unnecessary copies
type customConvert struct{}

func (*customConvert) stringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func (*customConvert) bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))

}

//SupportCustomConversion used for whether the user-defined conversion is supported. If so, the program performance will be improved
//customConvert implements zero copy. It depends on the underlying slice and string. Therefore, it does not provide a stable interface. If the slice or string source code of the underlying go changes, it must be switched to the raw implementatio
func SupportCustomConversion(support bool) opt {
	return func(o *option) {
		switch support {
		case true:
			o.convert = &customConvert{}
		default:
			o.convert = &rawConvert{}
		}
	}
}

//Bytes used for conversion from bytes to string
func Bytes(s string) []byte {
	return options.convert.stringToBytes(s)
}

////String used for conversion from string to bytes
func String(b []byte) string {
	return options.convert.bytesToString(b)
}

package slice

import (
	"reflect"
	"unsafe"
)

type StringHeader struct {
	Data uintptr
	Len int
}

type SliceHeader struct {
	Data uintptr
	Len int
	Cap int
}

//string to slice
func stringToBytes(s string) (r []byte){
	stringHeader := *(*reflect.StringHeader)(unsafe.Pointer(&s))

	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len: stringHeader.Len,
		Cap:stringHeader.Len,
	}

	r = *(*[]byte)(unsafe.Pointer(&bh))
	return
}

//slice to string
func sliceToString(b []interface{})(r string){
	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len: sliceHeader.Len,
	}

	r = *(*string)(unsafe.Pointer(&sh))
	return
}

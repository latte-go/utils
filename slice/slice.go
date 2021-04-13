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

func Max(vals []int32) int32 {
	var max int32
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func Min(vals []int32) int32{
	var min int32
	for _, val := range vals {
		if min == 0 || val <= min {
			min = val
		}
	}
	return min
}

func RemoveSlice(slc []int64)[]int64{
	result := []int64{}  // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
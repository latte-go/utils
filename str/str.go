package str

import "strconv"

func StrToInt64(inp string) (r int64) {
	var err error
	if r, err = strconv.ParseInt(inp, 10, 64); err != nil {
		return
	}
	return
}

func StrToFloat(inp string) (r float64) {
	var err error
	if r, err = strconv.ParseFloat(inp, 10); err != nil {
		return
	}
	return
}

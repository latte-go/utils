package str

import "strconv"

func StrToInt64(inp string) (r int64) {
	var err error
	if r, err = strconv.ParseInt(inp, 10, 64); err != nil {
		return
	}
	return
}

package ip

import "github.com/thinkeridea/go-extend/exnet"

func IpToInt(ip string) (r uint) {
	var err error
	if r, err = exnet.IPString2Long(ip); err != nil {
		r = 0
	}
	return
}

func IpIntToString(ip uint) (r string) {
	var err error
	if r, err = exnet.Long2IPString(ip); err != nil {
		r = ""
	}
	return
}

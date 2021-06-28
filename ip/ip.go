package ip

import "github.com/thinkeridea/go-extend/exnet"

func IpToInt(ip string) (r uint) {
	var err error
	if r, err = exnet.IPString2Long(ip); err != nil {
		r = 0
	}
	return
}

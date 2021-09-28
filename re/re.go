package re

import (
	"regexp"
	"strings"
)

// todo 配合https://rubular.com/ 这个网站使用

func RegexGroup(regex, response string) (r []map[string]string) {

	r = make([]map[string]string, 0)
	// regex :=
	//解析字符串
	//替换正则为go 匹配的格式
	d := strings.ReplaceAll(regex, "?<", "?P<")

	res := regexp.MustCompile(d)
	matchs := res.FindAllStringSubmatch(response, -1)
	groupNames := res.SubexpNames()

	for _, match := range matchs {
		temp := make(map[string]string, 0)
		if len(match) != len(groupNames) {
			continue
		}

		for k, v := range groupNames {
			if k == 0 && len(v) == 0 {
				continue
			}
			temp[v] = match[k]
		}
		r = append(r, temp)
	}
	return
}

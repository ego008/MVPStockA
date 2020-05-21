package file

import (
	"net/url"
	"strings"
)

func Report(code []string) {
	pushUrl := url.Values{}
	pushUrl.Set("data", strings.Join(code, ","))
	go Get("http://push.strcpy.cn/gupiao/push", pushUrl, nil, 5000)
}

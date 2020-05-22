package file

import "net/http"

func Report(data string) {
	h := http.Header{}
	go Post("http://push.strcpy.cn/gupiao/push", data, h, 5000)
}

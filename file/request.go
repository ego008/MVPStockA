package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func Get(requestUrl string, query url.Values, header http.Header, timeOut int64) string {
	if query.Encode() != "" {
		requestUrl = requestUrl + "?" + query.Encode()
	}
	u, err := url.Parse(requestUrl)
	if err != nil {
		return ""
	}
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Host:   u.Host,
		Header: header,
	}
	c := &http.Client{
		Timeout: time.Millisecond * time.Duration(timeOut),
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(data)
}

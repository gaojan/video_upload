package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type RequestInfo struct {
	Url           string
	Data          map[string]string //post要传输的数据，必须key value必须都是string
	DataInterface map[string]interface{}
}

//用于 application/x-www-form-urlencoded
func (r *RequestInfo) PostUrlEncoded() ([]byte, error) {
	client := &http.Client{}
	//post要提交的数据
	DataUrlVal := url.Values{}
	for key, val := range r.Data {
		DataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest("POST", r.Url, strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		return nil, err
	}
	//header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//提交请求
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	//读取返回值
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

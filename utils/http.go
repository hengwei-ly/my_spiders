package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

//获取网页内容 GET方式
func GetForContents(url string) ([]byte, error) {
	client := &http.Client{}
	request, req_err := http.NewRequest("GET", url, nil)
	if req_err != nil {
		return nil, req_err
	}

	request.Header.Set("Content-Type", "text/html")
	request.Header.Set("user-agent", `Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36`)
	resp, client_err := client.Do(request)
	if client_err != nil {
		return nil, client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		return nil, resp_err
	}
	return body, nil
}

//获取网页内容 POST方式
func PostForContents(url string, params map[string]interface{}) ([]byte, error) {
	client := &http.Client{}
	paramBs, err := json.Marshal(params)
	if err != nil {
		return nil, errors.New("序列化参数失败： " + err.Error())
	}
	request, req_err := http.NewRequest("POST", url, bytes.NewReader(paramBs))
	if req_err != nil {
		return nil, req_err
	}
	request.Header.Set("Content-Type", "text/html")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return nil, client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		return nil, resp_err
	}
	return body, nil
}

//参数更多  但是更灵活
func RequestForHttpContents(url, method string, params map[string]interface{},
	headers map[string]string) ([]byte, error) {
	client := &http.Client{}

	pa := &bytes.Reader{}
	if strings.ToLower(method) == "post" {
		paramBs, err := json.Marshal(params)
		if err != nil {
			return nil, errors.New("序列化参数失败： " + err.Error())
		}
		pa = bytes.NewReader(paramBs)
	} else {
		pa = nil
	}

	request, req_err := http.NewRequest(method, url, pa)
	if req_err != nil {
		return nil, req_err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	resp, client_err := client.Do(request)
	if client_err != nil {
		return nil, client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		return nil, resp_err
	}
	return body, nil
}

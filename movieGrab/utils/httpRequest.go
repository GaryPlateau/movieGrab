package utils

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 设置header信息
// param addHeaders map[string]string
// return headers map[string]string
func SetHtmlHeader(addHeaders map[string]string) (headers map[string]string) {
	headers = make(map[string]string)
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	headers["Accept-Charset"] = "utf-8"
	headers["Accept-Language"] = "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"
	headers["Accept-Encoding"] = "gzip, deflate, br"
	headers["Cache-Control"] = "no-cache"
	//headers["Connection"] = "keep-alive"
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	if addHeaders != nil {
		for key, val := range addHeaders {
			headers[key] = val
		}
	}
	return headers
}

// get请求
// param url string,
// param headers map[string]string,
// param datas map[string]string
// return *http.Response
func GetHttpRequest(reqUrl string, headers map[string]string, datas map[string]string, isProxy bool) *http.Response {
	var tmpParams, getReqUrl string

	getReqUrl = reqUrl
	if len(datas) != 0 {
		for key, val := range datas {
			tmpParams += key + "=" + val + "&"
		}
		params := tmpParams[:len(tmpParams)-1]
		getReqUrl = reqUrl + "?" + params
	}

	req, err := http.NewRequest("GET", getReqUrl, nil)
	if err != nil {
		fmt.Println("GET请求失败", err)
		return nil
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	if isProxy {
		proxy := "http://127.0.0.1:8080"
		proxyAddress, _ := url.Parse(proxy)
		timeout := 10 * time.Second
		client = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddress),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("GET客户端建立失败", err)
		return nil
	}

	if 200 == response.StatusCode {
		return response
	} else if 400 == response.StatusCode {
		fmt.Println("服务器请求失败", err)
		return nil
	} else {
		return nil
	}
}

func PostHttpRequest(requireUrl string, headers map[string]string, datas map[string]string, isProxy bool) []byte {
	var reader io.ReadCloser
	var requestBody string
	for key, val := range datas {
		requestBody += key + `=` + val + `&`
	}
	requestBody = requestBody[:len(requestBody)-1]
	stringReader := strings.NewReader(requestBody)

	request, err := http.NewRequest("POST", requireUrl, stringReader)
	if err != nil {
		fmt.Println("POST请求失败", err)
		return nil
	}

	for key, val := range headers {
		request.Header.Add(key, val)
	}

	//不认证ssl
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	if isProxy {
		proxy := "http://127.0.0.1:8080"
		proxyAddress, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddress),
			},
		}
	}

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("POST客户端建立失败", err)
		return nil
	}

	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else {
		reader = response.Body
	}

	if 200 == response.StatusCode {
		respBytes, _ := io.ReadAll(reader)
		return respBytes
	} else if 403 == response.StatusCode {
		fmt.Println("服务器拒绝响应")
	} else if 404 == response.StatusCode {
		fmt.Println("服务器未找到")
	}

	return nil
}

// 编码检测
// func determineEncoding(r io.Reader) encoding.Encoding {
// 	bytes, err := bufio.NewReader(r).Peek(1024)
// 	if err != nil {
// 		panic(err)
// 	}
// 	e, _, _ := charset.DetermineEncoding(bytes, "")
// 	return e
// }

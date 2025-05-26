package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"log"
	"net/http"
)

// 检查当前网站页面的编码
// param reqUrl string
// return code encoding.Encoding
func CheckWebDecode(reqUrl string) (code encoding.Encoding) {
	response, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println("get请求失败", err)
		return nil
	}
	defer response.Body.Close()
	//headers := api.GetDyttHeader()
	//response := GetHttpRequest(reqUrl, headers, nil, true)

	body := bufio.NewReader(response.Body)
	bytes, err := body.Peek(1024 * 4) //读取1024个字节进行判断其编码
	if err != nil {
		log.Printf("读取网站内容失败:%v", err)
		return nil
	}

	code, _, _ = charset.DetermineEncoding(bytes, "") //读取1024个字节进行判断其编码
	return code
}

// 通过编码解析文本
// param code encoding.Encoding, content string
// return result string
func DecodeAnyCode(code encoding.Encoding, content string) (result string) {
	//读取并打印获取的信息
	result, _, err := transform.String(code.NewDecoder(), content)
	if err != nil {
		fmt.Println("编译内容失败", err)
		return ""
	}
	return
}

func DecodeGBK(tByte []byte) ([]byte, error) {
	// req, _ := http.NewRequest("GET", movieUrl, nil)
	// req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	// req.Header.Add("Accept-Charset", "utf-8")
	// req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	// req.Header.Add("Cache-Control", "no-cache")
	// req.Header.Add("Connection", "keep-alive")
	// client := &http.Client{}
	// response2, err := client.Do(req)
	// body2 := bufio.NewReader(response2.Body)
	// bytes2, err := body2.Peek(1024) //读取1024个字节进行判断其编码
	// if err != nil {
	// 	log.Printf("!!!!!!!!fetch error :%v", err)
	// 	return ""
	// }

	// e2, _, _ := charset.DetermineEncoding(bytes2, "") //读取1024个字节进行判断其编码

	// utf8Reader2 := transform.NewReader(body2, e2.NewDecoder())
	// //读取并打印获取的信息
	// result2, err := ioutil.ReadAll(utf8Reader2)
	// if err != nil {
	// 	fmt.Println("获取http内容失败", err)
	// 	return ""
	// }
	// fmt.Println(string(result2))
	// return string(result2)
	return nil, nil
}

// use iconv.Open or iconv.ConverString is bad
func oldFunc() {
	//defer response.Body.Close()

	//html := mahonia.NewDecoder("gbk").ConvertString(body)
	//fmt.Println(html)
	//doc.Find("#dede_content").Each(func(i int, s *goquery.Selection) {
	//html := s.Find("p").Text()
	//fmt.Println(html)
	// cd, err := iconv2.Open("utf-8", "gb2312")
	// if err != nil {
	// 	fmt.Println("charset:", charset, err)
	// }
	// defer cd.Close()
	// moiveInfo = cd.ConvString(html)
	// newHtml, err := iconv.ConvertString(html, charset, "utf-8")
	// if err != nil {
	// 	newHtml, err = iconv.ConvertString(html, "gb2312", "utf-8")
	// 	if err != nil {
	// 		newHtml, err = iconv.ConvertString(html, "GBK", "utf-8")
	// 		if err != nil {
	// 			cd, err := iconv2.Open("utf-8", "gb2312")
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 			defer cd.Close()
	// 			newHtml = cd.ConvString(html)
	// 			if newHtml == "" {
	// 				fmt.Println(charset)
	// 				fmt.Println(movieUrl)
	// 			}
	// 		}
	// 	}
	// }

	//})
	//fmt.Println(string(newHtml))
	// moiveInfo = string(newHtml)
	// return moiveInfo
}

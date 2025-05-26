package api

import "movieGrab/utils"

func GetAllMovieSearchUrl() (allMovieApi []string) {
	dygangPostUrl := "https://www.dygang.tv/e/search/index.php"
	ygdy8PostUrl := "http://s.ygdy8.com/plus/s01.php?typeid=1&keyword=%BF%F1%EC%AD"
	ys66PostUrl := "https://www.66ys.cc/e/search/index.php"
	allMovieApi = append(allMovieApi, dygangPostUrl)
	allMovieApi = append(allMovieApi, ygdy8PostUrl)
	allMovieApi = append(allMovieApi, ys66PostUrl)
	return
}

func GetDygangHeader() map[string]string {
	var headers map[string]string
	headers = make(map[string]string, 20)
	headers["Accept-Encoding"] = "gzip, deflate"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Origin"] = "http://www.dygang.tv"
	headers["Referer"] = "http://www.dygang.tv"
	headers["Sec-Fetch-Dest"] = "document"
	headers["Sec-Fetch-Mode"] = "navigate"
	headers["Sec-Fetch-Site"] = "cross-site"
	headers["Sec-Fetch-User"] = "?1"
	headers["Upgrade-Insecure-Requests"] = "1"

	headers = utils.SetHtmlHeader(headers)
	return headers
}

func GetDyttHeader() map[string]string {
	var headers map[string]string
	headers = make(map[string]string, 20)
	headers["Host"] = "www.dyttcn.com"
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0"
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	headers["Accept-Language"] = "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"
	headers["Accept-Encoding"] = "gzip, deflate, br"
	headers["Upgrade-Insecure-Requests"] = "1"
	headers["Sec-Fetch-Dest"] = "document"
	headers["Sec-Fetch-Mode"] = "navigate"
	headers["Sec-Fetch-Site"] = "cross-site"
	headers["Priority"] = "u=1"
	headers["Pragma"] = "no-cache"
	headers["Cache-Control"] = "no-cache"
	headers["Te"] = "trailers"
	headers["Connection"] = "close"

	headers = utils.SetHtmlHeader(headers)
	return headers
}

func Get66ysHeader() map[string]string {
	var headers map[string]string
	headers = make(map[string]string, 20)
	headers["Host"] = "www.66yingshi.com"
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	headers["Accept-Encoding"] = "gzip, deflate, br"
	headers["Accept-Language"] = "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Sec-Fetch-Dest"] = "document"
	headers["Sec-Fetch-Mode"] = "navigate"
	headers["Sec-Fetch-Site"] = "cross-site"
	headers["Sec-Fetch-User"] = "?1"
	headers["Priority"] = "u=1"
	headers["Te"] = "trailers"
	headers["Connection"] = "close"
	headers["Upgrade-Insecure-Requests"] = "1"

	headers = utils.SetHtmlHeader(headers)
	return headers
}

func GetDygangUrl() (basicUrl string, searchUrl string) {
	basicUrl = "https://www.dygang.tv"
	searchUrl = "https://www.dygang.tv/e/search/result/"
	return
}

func GetDygangMoiveKind(moiveKind string) map[string]string {
	var datas map[string]string
	datas = make(map[string]string)
	switch moiveKind {
	case "kongbupian":
		{
			datas["searchid"] = "178275"
		}
	case "xijupian":
		{
			datas["searchid"] = "176303"
		}
	case "dongzuopian":
		{
			datas["searchid"] = "178216"
		}
	case "aiqingpian":
		{
			datas["searchid"] = "176518"
		}
	case "kehuanpian":
		{
			datas["searchid"] = "184707"
		}
	case "zhanzhengpian":
		{
			datas["searchid"] = "178108"
		}
	case "xuanyipian":
		{
			datas["searchid"] = "179239"
		}
	default:
		{
			datas["searchid"] = ""
		}
	}
	return datas
}

func GetDytt8Url() (basicUrl string, searchUrl string) {
	basicUrl = "https://www.dyttcn.com"
	searchUrl = "https://www.dyttcn.com/plus/search.php"
	return
}

func Get66ysUrl() (basicUrl string, searchUrl string) {
	basicUrl = "https://www.5266ys.com/"
	searchUrl = "https://www.66yingshi.com"
	return
}

func Get66ysMoiveKind(basicUrl string, moiveKind string) {
	switch moiveKind {
	case "kongbupian":
		{
			basicUrl = basicUrl + "kongbupian/"
		}
	case "xijupian":
		{
			basicUrl = basicUrl + "xijupian/"
		}
	case "dongzuopian":
		{
			basicUrl = basicUrl + "dongzuopian/"
		}
	case "aiqingpian":
		{
			basicUrl = basicUrl + "aiqingpian/"
		}
	case "kehuanpian":
		{
			basicUrl = basicUrl + "kehuanpian/"
		}
	case "zhanzhengpian":
		{
			basicUrl = basicUrl + "zhanzhengpian/"
		}
	case "jilupian":
		{
			basicUrl = basicUrl + "jilupian/"
		}
	case "juqingpian":
		{
			basicUrl = basicUrl + "juqingpian/"
		}
	case "guochanju":
		{
			basicUrl = basicUrl + "dsj/"
		}
	case "gangtaiju":
		{
			basicUrl = basicUrl + "dsj2/"
		}
	case "rihanju":
		{
			basicUrl = basicUrl + "dsj2/"
		}
	case "oumeiju":
		{
			basicUrl = basicUrl + "dsj2/"
		}
	case "guopeidianying":
		{
			basicUrl = basicUrl + "gy/"
		}
	default:
		{
			basicUrl = basicUrl + ""
		}
	}
}

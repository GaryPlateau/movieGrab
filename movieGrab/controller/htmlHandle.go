package controller

import (
	"fmt"
	"maps"
	"movieGrab/api"
	"movieGrab/config"
	"movieGrab/utils"
	"net/http"
	"os"
)

type HtmlHandler struct {
	basicUrl  string
	searchUrl string
	moiveKind string
	conf      *config.ConfigInfo
}

func (this *HtmlHandler) SetConfigInfo() {
	this.conf = this.conf.LoadConfig()
}

func (this *HtmlHandler) SetMoiveKind(kind string) {
	this.moiveKind = kind
}

func (this *HtmlHandler) SetBasicUrl(webName string) {
	switch webName {
	case "dygang":
		this.basicUrl, this.searchUrl = api.GetDygangUrl()
	case "dytt8":
		this.basicUrl, this.searchUrl = api.GetDytt8Url()
	case "66ys":
		this.basicUrl, this.searchUrl = api.Get66ysUrl()
	default:
		os.Exit(404)
	}
}

// 通过url返回对应doc对象
// param url string
// param datas map[string]string
// return doc, resp
func (this *HtmlHandler) initRequest(url string, setHeader, datas map[string]string) (response *http.Response) {
	var headers map[string]string
	headers = utils.SetHtmlHeader(headers)
	if setHeader != nil {
		maps.Copy(headers, setHeader)
	}
	response = utils.GetHttpRequest(url, headers, datas, true)

	if response == nil {
		return
	}
	if 200 != response.StatusCode {
		fmt.Println("网页请求失败，状态:", response.Status)
		return nil
	}
	return
}

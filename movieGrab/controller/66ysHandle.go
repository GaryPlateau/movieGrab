package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/transform"
	"io"
	"movieGrab/api"
	"movieGrab/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Ys66Handler struct {
	movieKind     string
	htmlHandler   *HtmlHandler
	processStatus string
}

func NewYs66Handler() (dh *Ys66Handler) {
	dh = &Ys66Handler{
		"dyzz",
		new(HtmlHandler),
		"",
	}
	return
}

func (this *Ys66Handler) Ys66Menu(webSelect string) {
	switch webSelect {
	case "1":
		this.movieKind = "dongzuopian"
	case "2":
		this.movieKind = "kongbupian"
	case "3":
		this.movieKind = "zhanzhengpian"
	case "4":
		this.movieKind = "kehuanpian"
	case "5":
		this.movieKind = "aiqingpian"
	case "6":
		this.movieKind = "xijupian"
	case "7":
		this.movieKind = "jilupian"
	case "8":
		this.movieKind = "bd"
	case "9":
		this.movieKind = "dsj"
	case "10":
		this.movieKind = "dsj2"
	case "11":
		this.movieKind = "dsj3"
	case "12":
		this.movieKind = "dsj4"
	case "13":
		this.movieKind = "gy"
	default:
		this.movieKind = ""
	}
}

func (this *Ys66Handler) GetProcessStatus() string {
	return this.processStatus
}

func (this *Ys66Handler) Get66dyPages() (movieTotal int, pageTotal int) {
	var movieTotalChan chan string
	movieTotalChan = make(chan string, 1)
	this.htmlHandler.SetBasicUrl("66ys")
	this.htmlHandler.SetMoiveKind(this.movieKind)
	go this.getPageTotal(movieTotalChan)
	movieTotal, _ = strconv.Atoi(<-movieTotalChan)
	if movieTotal%20 == 0 {
		pageTotal = movieTotal / 20
	} else {
		pageTotal = movieTotal/20 + 1
	}
	if 0 == movieTotal {
		fmt.Println("网站加载超时请重试")
		//exitChan <- true
		return
	}
	return
}

func (this *Ys66Handler) GetMovieInfo(wg *sync.WaitGroup, startPage int, endPage int, pageTotal int, movieTotal int) (movieInfoMap map[int]map[string]string) {
	var pageUrlChan chan string
	var moiveInfoChan chan string
	this.htmlHandler.SetConfigInfo()
	moiveCount := 0
	if startPage == endPage {
		moiveCount = 20
	} else {
		if pageTotal == endPage {
			moiveCount = (endPage-1)*20 + movieTotal%20
		} else {
			moiveCount = endPage * 20
		}
	}
	pageUrlChan = make(chan string, moiveCount)

	for i := startPage; i <= endPage; i++ {
		wg.Add(1)
		go this.getPageMovieUrl(wg, pageUrlChan, strconv.Itoa(i))
		time.Sleep(time.Millisecond * 50)
	}
	wg.Wait()

	for {
		if len(pageUrlChan) == moiveCount {
			close(pageUrlChan)
			break
		}
	}

	pageUrlChanLen := len(pageUrlChan)
	moiveInfoChan = make(chan string, pageUrlChanLen)
	wg.Add(pageUrlChanLen)
	count := 0
	for moiveUrl := range pageUrlChan {
		count++
		go this.getMoiveInfo(wg, moiveUrl, moiveInfoChan)
		time.Sleep(time.Millisecond * 50)
		this.processStatus = strconv.FormatFloat((float64(count)/float64(pageUrlChanLen))*100, 'f', 2, 64) + "%"
	}
	wg.Wait()
	close(moiveInfoChan)

	movieInfoMap = make(map[int]map[string]string, len(moiveInfoChan))
	count = 0
	for info := range moiveInfoChan {
		movieInfoMap[count] = utils.MovieInfo66ysHander(info)
		count++
	}
	return movieInfoMap
}

func (this *Ys66Handler) getPageTotal(movieTotalChan chan string) {
	response := this.htmlHandler.initRequest(this.htmlHandler.basicUrl+this.htmlHandler.moiveKind+"/index.html", nil, nil)
	if response == nil {
		movieTotalChan <- "0"
		close(movieTotalChan)
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}
	//div class="x"
	doc.Find("a[title='Total record']").Each(func(i int, s *goquery.Selection) {
		pageTotal := s.Text()
		reg := `\d+`
		re := regexp.MustCompile(reg)
		res := re.FindStringSubmatch(pageTotal)
		movieTotalChan <- res[0]
	})
	close(movieTotalChan)
}

func (this *Ys66Handler) getPageMovieUrl(wg *sync.WaitGroup, pageUrlChan chan string, nowPage string) {
	defer wg.Done()
	var response *http.Response
	headers := api.Get66ysHeader()
	if nowPage == "1" {
		response = this.htmlHandler.initRequest(this.htmlHandler.searchUrl+"/"+this.htmlHandler.moiveKind+"/index.html", headers, nil)
	} else {
		response = this.htmlHandler.initRequest(this.htmlHandler.searchUrl+"/"+this.htmlHandler.moiveKind+"/index_"+nowPage+".html", headers, nil)
	}
	if response == nil {
		pageUrlChan <- ""
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}

	doc.Find(".listimg").Each(func(i int, s *goquery.Selection) {
		aUrl, _ := s.Find("a").Attr("href")
		pageUrlChan <- aUrl
	})
	response.Body.Close()
}

func (this *Ys66Handler) getMoiveInfo(wg *sync.WaitGroup, moiveUrl string, moiveInfoChan chan string) {
	defer wg.Done()
	failExportPath := this.htmlHandler.conf.FailExportPath
	failExportPath = failExportPath + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	moiveInfo := this.getMoiveDetails(moiveUrl)
	moiveInfo = strings.Replace(moiveInfo, "'", "", -1)
	if moiveInfo == "" {
		utils.WriteFile(failExportPath, []byte(moiveInfo+"\n"))
	}
	moiveInfoChan <- moiveInfo
}

func (this *Ys66Handler) getMoiveDetails(movieUrl string) (moiveInfo string) {
	headers := api.Get66ysHeader()
	code := utils.CheckWebDecode(this.htmlHandler.searchUrl + movieUrl)
	response := this.htmlHandler.initRequest(this.htmlHandler.searchUrl+movieUrl, headers, nil)
	if response == nil {
		return
	}
	defer response.Body.Close()

	body := bufio.NewReader(response.Body)
	utf8Reader := transform.NewReader(body, code.NewDecoder())

	result, err := io.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println("获取http内容失败", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(result))
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}

	doc.Find("#text").Each(func(i int, s *goquery.Selection) {
		basicInfo := s.Text()
		moiveInfo += basicInfo
	})

	moiveInfo += "◎电影链接　"
	doc.Find("td[bgcolor='#ffffbb']").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		url, _ := s.Find("a").Attr("href")
		moiveInfo += name + "*:" + url + "\r\n"
	})
	return
}

func (this *Ys66Handler) ExportDataToExcel(movieInfoMap map[int]map[string]string) bool {
	modelPath := this.htmlHandler.conf.ModelPath + "movie.xlsx"
	savePath := this.htmlHandler.conf.SavePath
	savePath += this.movieKind + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	eh := utils.CreateExcelFile(modelPath, savePath)
	eh.CreateExcelTitle("movie1")
	err := eh.OpenExcelFile("movie1")
	if err != nil {
		return false
	} else {
		for _, infos := range movieInfoMap {
			eh.WriteContentToExcel(infos)
		}
		eh.SaveAsExcel()
		eh.CloseExcelFile()
		return true
	}
}

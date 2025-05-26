package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/djimenez/iconv-go"
	"golang.org/x/text/transform"
	"io/ioutil"
	"movieGrab/api"
	"movieGrab/utils"
	"strconv"
	"sync"
	"time"
)

type DygangHandler struct {
	movieKind     string
	htmlHandler   *HtmlHandler
	processStatus string
}

// // 判断一个页面自有的charset设置
// // param url string
// // return charset string, doc *goquery.Document
//
//	func (this *DygangHandler) getWebCharset(url string) (charset string, doc *goquery.Document) {
//		response := this.htmlHandler.initRequest(this.htmlHandler.searchUrl, nil)
//		defer response.Body.Close()
//
//		doc, err := goquery.NewDocumentFromReader(response.Body)
//		if err != nil {
//			fmt.Println("获取doc节点错误", err)
//			return
//		}
//
//		pattern := `(.+)(\s)(charset=)(.*)`
//		re := regexp.MustCompile(pattern)
//
//		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
//			char, flag := s.Attr("content")
//			if flag {
//				matchArr := re.FindStringSubmatch(char)
//				if len(matchArr) != 0 {
//					charset = matchArr[4]
//				}
//			}
//		})
//		return charset, doc
//	}

func NewDygangHandler() (dh *DygangHandler) {
	dh = &DygangHandler{
		"all",
		new(HtmlHandler),
		"",
	}
	return
}

func (this *DygangHandler) DygangMenu(webSelect string) {
	switch webSelect {
	case "1":
		this.movieKind = "kongbupian"
	case "2":
		this.movieKind = "xijupian"
	case "3":
		this.movieKind = "dongzuopian"
	case "4":
		this.movieKind = "aiqingpian"
	case "5":
		this.movieKind = "kehuanpian"
	case "6":
		this.movieKind = "zhanzhengpian"
	case "7":
		this.movieKind = "xuanyipian"
	default:
		this.movieKind = "all"
	}
}

func (this *DygangHandler) GetProcessStatus() string {
	return this.processStatus
}

func (this *DygangHandler) GetDygangPages() (movieTotal int, pageTotal int) {
	var pageTotalChan chan string
	pageTotalChan = make(chan string, 1)
	this.htmlHandler.SetBasicUrl("dygang")
	this.htmlHandler.SetMoiveKind(this.movieKind)
	go this.getPageTotal(pageTotalChan)
	movieTotal, _ = strconv.Atoi(<-pageTotalChan)
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

func (this *DygangHandler) GetMovieInfo(wg *sync.WaitGroup, startPage int, endPage int, pageTotal int, movieTotal int) (movieInfoMap map[int]map[string]string) {
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
		go this.getPageMovieUrl(wg, i, pageUrlChan)
		time.Sleep(time.Second)
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
		time.Sleep(time.Millisecond * 10)
		this.processStatus = strconv.FormatFloat((float64(count)/float64(pageUrlChanLen))*100, 'f', 2, 64) + "%"
	}
	wg.Wait()
	close(moiveInfoChan)

	movieInfoMap = make(map[int]map[string]string, len(moiveInfoChan))
	count = 0
	for info := range moiveInfoChan {
		movieInfoMap[count] = utils.MovieInfoHander(info)
		count++
	}
	return movieInfoMap
}

func (this *DygangHandler) ExportDataToExcel(movieInfoMap map[int]map[string]string) bool {
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
		//eh.ReadOrginContents()
		//eh.AppendContents(movieInfoMap)
		return true
	}

	// for index, movieInfo := range movieInfoMap {
	// 	wg.Add(1)
	// 	utils.WriteExcel(&wg, "movie1", movieInfo, index+1)
	// }
	// filePath := "C:/Users/Administrator/Desktop/movie.xlsx"
	// //utils.SaveExcelFile
	// wg.Wait()

	//for {
	//	if 0 == len(pageUrlChan) {
	//		exitChan <- true
	//		break
	//	}
	//}
}

// 获取每页内容
func (this *DygangHandler) getMoiveInfo(wg *sync.WaitGroup, moiveUrl string, moiveInfoChan chan string) {
	defer wg.Done()
	failExportPath := this.htmlHandler.conf.FailExportPath
	failExportPath = failExportPath + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	moiveInfo := this.getMoiveDetails(moiveUrl)
	if moiveInfo == "" {
		utils.WriteFile(failExportPath, []byte(moiveInfo+"\n"))
	}
	moiveInfoChan <- moiveInfo
}

// 获取电影类型总页数
func (this *DygangHandler) getPageTotal(pageTotalChan chan string) {
	datas := api.GetDygangMoiveKind(this.htmlHandler.moiveKind)
	response := this.htmlHandler.initRequest(this.htmlHandler.searchUrl, nil, datas)
	if response == nil {
		pageTotalChan <- "0"
		close(pageTotalChan)
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}

	doc.Find("a[title]").Each(func(i int, s *goquery.Selection) {
		pageTotal := s.Find("b").Text()
		pageTotal, err := iconv.ConvertString(pageTotal, "gb2312", "utf-8")
		if err != nil {
			fmt.Println("转换编码失败", err)
			return
		}
		pageTotalChan <- pageTotal
	})
	close(pageTotalChan)
}

// 获取每页body的url
// param page int
// parma pageUrlChan chan string
func (this *DygangHandler) getPageMovieUrl(wg *sync.WaitGroup, page int, pageUrlChan chan string) {
	defer wg.Done()
	datas := api.GetDygangMoiveKind(this.htmlHandler.moiveKind)
	datas["page"] = strconv.Itoa(page - 1)
	response := this.htmlHandler.initRequest(this.htmlHandler.searchUrl, nil, datas)
	if response == nil {
		pageUrlChan <- ""
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}

	doc.Find(".border1").Each(func(i int, s *goquery.Selection) {
		hrefText, errBool := s.Find("a").Attr("href")
		if errBool != true {
			fmt.Println("不存在此属性", i)
		}
		//fmt.Println(fmt.Sprintf("%v%v", this.basicUrl, hrefText))
		pageUrlChan <- fmt.Sprintf("%v%v", this.htmlHandler.basicUrl, hrefText)
	})
	//fmt.Println("读取电影页面数量:" + strconv.Itoa(len(pageUrlChan)))
	response.Body.Close()
}

// 获取电影具体详情
// param movieUrl string
// return moiveInfo string
func (this *DygangHandler) getMoiveDetails(movieUrl string) (moiveInfo string) {
	code := utils.CheckWebDecode(movieUrl)
	response := this.htmlHandler.initRequest(movieUrl, nil, nil)
	if response == nil || code == nil {
		return
	}
	defer response.Body.Close()

	body := bufio.NewReader(response.Body)
	utf8Reader := transform.NewReader(body, code.NewDecoder())

	result, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println("获取http内容失败", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(result))
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}
	doc.Find("#dede_content").Each(func(i int, s *goquery.Selection) {
		moiveInfo = s.Find("p").Text()
	})
	movieSeed := "◎电影链接  "
	doc.Find("td[bgcolor='#ffffbb']").Each(func(i int, s *goquery.Selection) {
		seed, _ := s.Find("a").Attr("href")
		//movieSeed = movieSeed + strconv.Itoa(i) + ":" + mahonia.NewDecoder("gbk").ConvertString(seed) + "  "
		movieSeed = movieSeed + "*:" + utils.DecodeAnyCode(code, seed) + "  "
	})

	moiveInfo = moiveInfo + "\n" + movieSeed + "\n"
	return
}

// 获取电影种子连接
// param movieUrl string
// return movieSeed string
func (this *DygangHandler) getMoiveDownloadUrl(movieUrl string) (movieSeed string) {
	response := this.htmlHandler.initRequest(movieUrl, nil, nil)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("获取doc节点错误", err)
		return
	}

	//length := doc.Find("td[bgcolor='#ffffbb']").Length()
	doc.Find("td[bgcolor='#ffffbb']").Each(func(i int, s *goquery.Selection) {
		seed, _ := s.Find("a").Attr("href")
		movieSeed = movieSeed + strconv.Itoa(i) + ":" + mahonia.NewDecoder("gbk").ConvertString(seed) + " "
	})

	return movieSeed
}

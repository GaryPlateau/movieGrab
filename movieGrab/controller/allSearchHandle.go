package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/gogf/gf/encoding/gcharset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"movieGrab/api"
	"movieGrab/utils"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type AllSearchMovieHandler struct {
	movieApi      []string
	htmlHandler   *HtmlHandler
	processStatus string
}

func NewAllSearchMovieHandler() (dh *AllSearchMovieHandler) {
	dh = &AllSearchMovieHandler{
		api.GetAllMovieSearchUrl(),
		new(HtmlHandler),
		"",
	}
	return
}

func (this *AllSearchMovieHandler) GetProcessStatus() string {
	return this.processStatus
}

func (this *AllSearchMovieHandler) GetAllSearchMovieUrl() []string {
	return this.movieApi
}

func (this *AllSearchMovieHandler) GetAllSearchResult(keyword string) string {
	tmpStr, err := gcharset.UTF8To("gbk", keyword)
	if err != nil {
		panic(err)
	}
	keywordGbk := ""
	for _, b := range []byte(tmpStr) {
		keywordGbk += fmt.Sprintf("%%%X", b)
	}
	dygangHeaders := api.GetDygangHeader()
	dygangDatas := make(map[string]string, 5)
	dygangDatas["tempid"] = "1"
	dygangDatas["tbname"] = "article"
	dygangDatas["keyboard"] = keywordGbk
	dygangDatas["show"] = "title%2Csmalltext"
	dygangDatas["submit"] = "%CB%D1%CB%F7"

	dygangPostResponse := utils.PostHttpRequest(this.movieApi[0], dygangHeaders, dygangDatas, false)
	pattern := `result\/\?searchid=(\d+)`
	re := regexp.MustCompile(pattern)
	address := re.FindStringSubmatch(string(dygangPostResponse))
	if len(address) == 0 {
		return ""
	} else {
		if address[1] == "0" {
			return ""
		} else {
			this.htmlHandler.basicUrl, _ = api.GetDygangUrl()
			this.htmlHandler.searchUrl = this.GetAllSearchMovieUrl()[0][:len(this.GetAllSearchMovieUrl()[0])-9] + "result/index.php"
			this.htmlHandler.moiveKind = address[1]
			return address[0]
		}
	}
}

func (this *AllSearchMovieHandler) GetMovieInfo(wg *sync.WaitGroup, startPage int, endPage int, pageTotal int, movieTotal int) (movieInfoMap map[int]map[string]string) {
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
	close(pageUrlChan)
	//for {
	//	if len(pageUrlChan) == moiveCount {
	//		close(pageUrlChan)
	//		break
	//	}
	//}

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

func (this *AllSearchMovieHandler) GetAllResultPages(resultUrl string) (movieTotal int, pageTotal int) {
	var movieTotalChan chan string
	movieTotalChan = make(chan string, 1)
	go this.getAllMovies(movieTotalChan, resultUrl)
	movieTotal, _ = strconv.Atoi(<-movieTotalChan)
	if 0 == movieTotal {
		fmt.Println("网站加载超时请重试")
		//exitChan <- true
		return
	}
	if movieTotal%20 == 0 {
		pageTotal = movieTotal / 20
	} else {
		pageTotal = movieTotal/20 + 1
	}
	return
}

func (this *AllSearchMovieHandler) getAllMovies(movieTotalChan chan string, resultUrl string) {
	response := this.htmlHandler.initRequest(resultUrl, nil, nil)
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

	doc.Find("a[title]").Each(func(i int, s *goquery.Selection) {
		pageTotal := s.Find("b").Text()
		if pageTotal == "" {
			return
		}
		pageTotal, err := iconv.ConvertString(pageTotal, "gb2312", "utf-8")
		if err != nil {
			fmt.Println("转换编码失败", err)
			return
		}
		movieTotalChan <- pageTotal
	})

	movieCount := 0
	doc.Find("td[width='250']").Each(func(i int, s *goquery.Selection) {
		movieUrl, _ := s.Find("a").Attr("href")
		if movieUrl != "" {
			movieCount++
		}
	})
	if len(movieTotalChan) == 0 {
		movieTotalChan <- strconv.Itoa(movieCount)
	}

	close(movieTotalChan)
}

// 获取每页body的url
// param page int
// parma pageUrlChan chan string
func (this *AllSearchMovieHandler) getPageMovieUrl(wg *sync.WaitGroup, page int, pageUrlChan chan string) {
	defer wg.Done()
	var datas map[string]string
	datas = make(map[string]string, 2)
	datas["searchid"] = this.htmlHandler.moiveKind
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

	doc.Find("td[width='250']").Each(func(i int, s *goquery.Selection) {
		movieUrl, _ := s.Find("a").Attr("href")
		if movieUrl != "" {
			pageUrlChan <- fmt.Sprintf("%v%v", this.htmlHandler.basicUrl, movieUrl)
		}
	})
	response.Body.Close()
}

// 获取每页内容
func (this *AllSearchMovieHandler) getMoiveInfo(wg *sync.WaitGroup, moiveUrl string, moiveInfoChan chan string) {
	defer wg.Done()
	failExportPath := this.htmlHandler.conf.FailExportPath
	failExportPath = failExportPath + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	moiveInfo := this.getMoiveDetails(moiveUrl)
	if moiveInfo == "" {
		utils.WriteFile(failExportPath, []byte(moiveInfo+"\n"))
	}
	moiveInfoChan <- moiveInfo
}

// 获取电影具体详情
// param movieUrl string
// return moiveInfo string
func (this *AllSearchMovieHandler) getMoiveDetails(movieUrl string) (moiveInfo string) {
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
		movieSeed = movieSeed + strconv.Itoa(i) + ":" + utils.DecodeAnyCode(code, seed) + "  "
	})

	moiveInfo = moiveInfo + "\n" + movieSeed + "\n"
	return
}

func (mw *AllSearchMovieHandler) ExportDataToExcel(movieInfoMap map[int]map[string]string) bool {
	modelPath := mw.htmlHandler.conf.ModelPath + "movie.xlsx"
	savePath := mw.htmlHandler.conf.SavePath
	savePath += mw.htmlHandler.moiveKind + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
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

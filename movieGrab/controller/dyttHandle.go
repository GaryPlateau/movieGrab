package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"golang.org/x/text/transform"
	"io"
	"movieGrab/utils"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Dytt8Handler struct {
	moiveKindID   string
	movieKind     string
	htmlHandler   *HtmlHandler
	processStatus string
}

func NewDytt8Handler() (dh *Dytt8Handler) {
	dh = &Dytt8Handler{
		"1",
		"dongzuopian",
		new(HtmlHandler),
		"",
	}
	return
}

func (this *Dytt8Handler) Dytt8Menu(webSelect string) {
	this.moiveKindID = webSelect
	switch webSelect {
	case "1":
		this.movieKind = "dongzuopian"
	case "2":
		this.movieKind = "juqingpian"
	case "3":
		this.movieKind = "aiqingpian"
	case "4":
		this.movieKind = "xijupian"
	case "5":
		this.movieKind = "kehuanpian"
	case "6":
		this.movieKind = "kongbupian"
	case "7":
		this.movieKind = "donghuapian"
	case "8":
		this.movieKind = "jingsongpian"
	case "9":
		this.movieKind = "zhanzhengpian"
	case "10":
		this.movieKind = "fanzuipian"
	case "11":
		this.movieKind = "zainanpian"
		this.moiveKindID = "19"
	case "12":
		this.movieKind = "jilupian"
		this.moiveKindID = "17"
	case "13":
		this.movieKind = "qihuanpian"
		this.moiveKindID = "14"
	default:
		this.movieKind = "all"
	}
	//case "1":
	//	this.movieKind = "dyzz"
	//case "2":
	//	this.movieKind = "china"
	//case "3":
	//	this.movieKind = "oumei"
	//case "4":
	//	this.movieKind = "hytv"
	//case "5":
	//	this.movieKind = "rihantv"
	//case "6":
	//	this.movieKind = "oumeitv"
	//case "7":
	//	this.movieKind = "zongyi2013"
	//case "8":
	//	this.movieKind = "dongman"
	//default:
	//	this.movieKind = "all"
	//}
}

func (this *Dytt8Handler) GetProcessStatus() string {
	return this.processStatus
}

func (this *Dytt8Handler) GetDytt8Pages() (movieTotal int, pageTotal int) {
	var pageTotalChan chan string
	pageTotalChan = make(chan string, 2)
	this.htmlHandler.SetBasicUrl("dytt8")
	this.htmlHandler.SetMoiveKind(this.movieKind)
	go this.getPageTotal(pageTotalChan)
	pageTotal, _ = strconv.Atoi(<-pageTotalChan)
	movieTotal, _ = strconv.Atoi(<-pageTotalChan)
	if 0 == movieTotal {
		fmt.Println("网站加载超时请重试")
		//exitChan <- true
		return
	}
	return
}

func (this *Dytt8Handler) GetMovieInfo(wg *sync.WaitGroup, startPage int, endPage int, pageTotal int, movieTotal int) (movieInfoMap map[int]map[string]string) {
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
		go this.getPageMovieUrl(wg, pageUrlChan, i)
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
		time.Sleep(time.Millisecond * 100)
		this.processStatus = strconv.FormatFloat((float64(count)/float64(pageUrlChanLen))*100, 'f', 2, 64) + "%"
	}
	wg.Wait()
	close(moiveInfoChan)

	movieInfoMap = make(map[int]map[string]string, len(moiveInfoChan))
	count = 0
	for info := range moiveInfoChan {
		movieInfoMap[count] = utils.MovieInfoDytt8Hander(info)
		count++
	}
	return movieInfoMap
}

// 获取电影类型总页数
func (this *Dytt8Handler) getPageTotal(pageTotalChan chan string) {
	//https://www.dyttcn.com/dongzuopian/list_1_1.html
	response := this.htmlHandler.initRequest(this.htmlHandler.basicUrl+"/"+this.htmlHandler.moiveKind+"/list_"+this.moiveKindID+"_1.html", nil, nil)
	if response == nil {
		pageTotalChan <- "0"
		close(pageTotalChan)
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("getPageTotal获取doc节点错误", err)
		return
	}
	//div class="x"
	doc.Find("span.pageinfo").Each(func(i int, s *goquery.Selection) {
		pageTotal := s.Text()
		pageTotal, err := iconv.ConvertString(pageTotal, "gb2312", "utf-8")
		if err != nil {
			fmt.Println("转换编码失败", err)
			return
		}
		reg := `共\s*(\d+)页(\d+)条`
		re := regexp.MustCompile(reg)
		res := re.FindStringSubmatch(pageTotal)
		pageTotalChan <- res[1]
		pageTotalChan <- res[2]
	})
	close(pageTotalChan)
}

func (this *Dytt8Handler) getPageMovieUrl(wg *sync.WaitGroup, pageUrlChan chan string, pageNo int) {
	defer wg.Done()
	response := this.htmlHandler.initRequest(this.htmlHandler.basicUrl+"/"+this.htmlHandler.moiveKind+"/list_"+this.moiveKindID+"_"+strconv.Itoa(pageNo)+".html", nil, nil)
	if response == nil {
		pageUrlChan <- ""
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("getPageMovieUrl获取doc节点错误", err)
		return
	}

	doc.Find("a[class=ulink]").Each(func(i int, s *goquery.Selection) {
		movieUrl, exists := s.Attr("href")
		if exists && i%2 == 0 {
			pageUrlChan <- fmt.Sprintf("%v%v", this.htmlHandler.basicUrl, movieUrl)
		}
	})
	response.Body.Close()
	return
}

func (this *Dytt8Handler) getMoiveInfo(wg *sync.WaitGroup, moiveUrl string, moiveInfoChan chan string) {
	defer wg.Done()
	failExportPath := this.htmlHandler.conf.FailExportPath
	failExportPath = failExportPath + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	moiveInfo := this.getMoiveDetails(moiveUrl)
	if moiveInfo == "" {
		utils.WriteFile(failExportPath, []byte(moiveInfo+"\n"))
	}
	moiveInfoChan <- moiveInfo
}

func (this *Dytt8Handler) getMoiveDetails(movieUrl string) (moiveInfo string) {
	code := utils.CheckWebDecode(movieUrl)
	response := this.htmlHandler.initRequest(movieUrl, nil, nil)
	if response == nil || code == nil {
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
		fmt.Println("getMoiveDetails获取doc节点错误", err)
		return
	}

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		moiveInfo += strings.Trim(s.Text(), "  ") + "\r\n"
	})
	movieSeed := "电影链接 "
	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		seed, _ := s.Find("a").Attr("href")
		movieSeed = movieSeed + utils.DecodeAnyCode(code, seed) + "  "
	})

	moiveInfo = moiveInfo + "\n" + movieSeed + "\n"
	return
}

func (this *Dytt8Handler) ExportDataToExcel(movieInfoMap map[int]map[string]string) bool {
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

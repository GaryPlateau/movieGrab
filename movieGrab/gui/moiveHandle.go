package gui

import (
	"github.com/lxn/walk"
	"log"
	"movieGrab/api"
	"movieGrab/controller"
	"regexp"
	"strconv"
	"time"
)

type MovieModel struct {
	walk.ListModelBase
	items []MovieItem
}

type MovieItems struct {
	movieItems []MovieItem
}

type MovieItem struct {
	name  string
	value string
}

func (mw *MyMainWindow) getMovieKinds() {
	if mw.dygMovie.Checked() {
		mw.cbMovieType.SetModel(api.GetDygangMovieType())
		mw.cbMovieType.SetCurrentIndex(-1)
	} else if mw.dyttMovie.Checked() {
		mw.cbMovieType.SetModel(api.GetDytt8MovieType())
		mw.cbMovieType.SetCurrentIndex(-1)
	} else if mw.dy66Movie.Checked() {
		mw.cbMovieType.SetModel(api.Get66ysMovieType())
		mw.cbMovieType.SetCurrentIndex(-1)
	} else if mw.allMovie.Checked() {
		mw.cbMovieType.SetModel(api.GetAllMovieType())
		mw.cbMovieType.SetCurrentIndex(-1)
	}
}

func (mw *MyMainWindow) getMoviePages(searchInfo *SearchInfo) {
	if err := mw.movieInfoDB.Submit(); err != nil {
		mw.statusTEHandle("统计失败,请稍后重试")
		walk.MsgBox(mw, "消息提示", "请选择完整信息", walk.MsgBoxIconInformation)
		log.Print(err)
		return
	}
	switch searchInfo.MovieWeb {
	case "dygang":
		dygHan = controller.NewDygangHandler()
		dygHan.DygangMenu(strconv.Itoa(searchInfo.MovieType))
		total, pages := dygHan.GetDygangPages()
		searchInfo.ResultTotalPages = pages
		searchInfo.ResultTotalMovies = total
		if total == 0 || pages == 0 {
			mw.statusTEHandle("统计失败,请稍后重试")
			return
		}
		mw.pagesLable.SetText("总共" + strconv.Itoa(total) + "部电影,共" + strconv.Itoa(pages) + "页，请选择页数：")
	case "dytt8":
		dytt8Han = controller.NewDytt8Handler()
		dytt8Han.Dytt8Menu(strconv.Itoa(searchInfo.MovieType))
		total, pages := dytt8Han.GetDytt8Pages()
		searchInfo.ResultTotalMovies = total
		searchInfo.ResultTotalPages = pages
		if total == 0 || pages == 0 {
			mw.statusTEHandle("统计失败,请稍后重试")
			return
		}
		mw.pagesLable.SetText("总共" + strconv.Itoa(total) + "部电影,共" + strconv.Itoa(pages) + "页，请选择页数：")
	case "66ys":
		ys66Han = controller.NewYs66Handler()
		ys66Han.Ys66Menu(strconv.Itoa(searchInfo.MovieType))
		total, pages := ys66Han.Get66dyPages()
		searchInfo.ResultTotalMovies = total
		searchInfo.ResultTotalPages = pages
		if total == 0 || pages == 0 {
			mw.statusTEHandle("统计失败,请稍后重试")
			return
		}
		mw.pagesLable.SetText("总共" + strconv.Itoa(total) + "部电影,共" + strconv.Itoa(pages) + "页，请选择页数：")
	default:

	}
	mw.statusTEHandle("统计完成")
	mw.contirm.SetEnabled(true)
	return
}

func (mw *MyMainWindow) getMovieInfo(searchInfo *SearchInfo) {
	var startPage, endPage, total, webPage int
	mw.statusSbiHandle("搜索中...", mw.failureIco)
	mw.statusTEHandle("搜索电影过程中,请稍后...")
	if "dygang" == searchInfo.MovieWeb || "66ys" == searchInfo.MovieWeb || "all" == searchInfo.MovieWeb {
		webPage = 20
	} else if "dytt8" == searchInfo.MovieWeb {
		webPage = 25
	}

	inputPage := mw.pagesNE.Text()
	reg := `(\d+)-(\d+)`
	re, _ := regexp.Compile(reg)
	found := re.FindStringSubmatch(inputPage)
	if len(found) > 0 {
		startPage, _ = strconv.Atoi(found[1])
		endPage, _ = strconv.Atoi(found[2])
		if endPage == searchInfo.ResultTotalPages {
			total = searchInfo.ResultTotalMovies
		} else {
			total = (endPage - startPage + 1) * webPage
		}
	} else {
		if endPage == searchInfo.ResultTotalPages {
			total = searchInfo.ResultTotalMovies - searchInfo.ResultTotalPages*webPage
		} else {
			total = webPage
		}
		startPage, _ = strconv.Atoi(inputPage)
		endPage = startPage
	}

	searchInfo.movieInfoMap = make(map[int]map[string]string, searchInfo.ResultTotalMovies)
	if mw.dygMovie.Checked() {
		go func() {
			for {
				time.Sleep(time.Second)
				mw.progressSbi.SetText("当前下载进度:" + dygHan.GetProcessStatus())
				if "100%" == dygHan.GetProcessStatus() {
					break
				}
			}
		}()
		searchInfo.movieInfoMap = dygHan.GetMovieInfo(&wg, startPage, endPage, searchInfo.ResultTotalPages, searchInfo.ResultTotalMovies)
	} else if mw.dyttMovie.Checked() {
		go func() {
			for {
				time.Sleep(time.Second)
				mw.progressSbi.SetText("当前下载进度:" + dytt8Han.GetProcessStatus())
				if "100%" == dytt8Han.GetProcessStatus() {
					break
				}
			}
		}()
		searchInfo.movieInfoMap = dytt8Han.GetMovieInfo(&wg, startPage, endPage, searchInfo.ResultTotalPages, searchInfo.ResultTotalMovies)
	} else if mw.dy66Movie.Checked() {
		go func() {
			for {
				time.Sleep(time.Second)
				mw.progressSbi.SetText("当前下载进度:" + ys66Han.GetProcessStatus())
				if "100%" == ys66Han.GetProcessStatus() {
					break
				}
			}
		}()
		searchInfo.movieInfoMap = ys66Han.GetMovieInfo(&wg, startPage, endPage, searchInfo.ResultTotalPages, searchInfo.ResultTotalMovies)
	} else if mw.allMovie.Checked() {
		go func() {
			for {
				time.Sleep(time.Second)
				mw.progressSbi.SetText("当前下载进度:" + allHan.GetProcessStatus())
				if "100%" == allHan.GetProcessStatus() {
					break
				}
			}
		}()
		searchInfo.movieInfoMap = allHan.GetMovieInfo(&wg, startPage, endPage, searchInfo.ResultTotalPages, searchInfo.ResultTotalMovies)
	}

	mw.statusSbiHandle("共搜索"+strconv.Itoa(total)+"部电影,"+"成功搜索到"+strconv.Itoa(len(searchInfo.movieInfoMap))+"部,丢失"+strconv.Itoa(total-len(searchInfo.movieInfoMap))+"部", mw.successIco)
	mw.export.SetEnabled(true)

}

func (mw *MyMainWindow) searchAllMovie(searchInfo *SearchInfo) {
	searchInfo.MovieWeb = "all"
	allHan = controller.NewAllSearchMovieHandler()
	searchResultUrl := allHan.GetAllSearchResult(mw.allMovieTE.Text())
	if searchResultUrl == "" {
		mw.statusTEHandle("没有找到关键词相关信息...")
		mw.statusSbiHandle("搜索失败,没有找到相关内容.", mw.failureIco)
		return
	}
	allSearchUrl := allHan.GetAllSearchMovieUrl()
	dyGangSearchResultUrl := allSearchUrl[0][:len(allSearchUrl[0])-9] + searchResultUrl
	searchInfo.ResultTotalMovies, searchInfo.ResultTotalPages = allHan.GetAllResultPages(dyGangSearchResultUrl)
	if searchInfo.ResultTotalMovies == 0 || searchInfo.ResultTotalPages == 0 {
		mw.statusTEHandle("统计失败,请稍后重试")
		mw.statusSbiHandle("统计失败,请稍后重试", mw.failureIco)
		return
	}
	mw.pagesLable.SetText("总共" + strconv.Itoa(searchInfo.ResultTotalMovies) + "部电影,共" + strconv.Itoa(searchInfo.ResultTotalPages) + "页，请选择页数：")
	mw.contirm.SetEnabled(true)
}

func (mw *MyMainWindow) exportMovieToExcel(searchInfo *SearchInfo) {
	var flag bool
	mw.statusSbiHandle("电影数据导出中...", mw.failureIco)
	if "dygang" == searchInfo.MovieWeb {
		flag = dygHan.ExportDataToExcel(searchInfo.movieInfoMap)
	} else if "dytt8" == searchInfo.MovieWeb {
		flag = dytt8Han.ExportDataToExcel(searchInfo.movieInfoMap)
	} else if "66ys" == searchInfo.MovieWeb {
		flag = ys66Han.ExportDataToExcel(searchInfo.movieInfoMap)
	} else if "all" == searchInfo.MovieWeb {
		flag = allHan.ExportDataToExcel(searchInfo.movieInfoMap)
	}

	if flag {
		walk.MsgBox(mw, "导出提示", "导出电影数据成功！", walk.MsgBoxUserIcon)
		mw.statusSbiHandle("电影数据导出成功!", mw.successIco)
	} else {
		walk.MsgBox(mw, "导出提示", "导出电影数据失败！", walk.MsgBoxIconError)
		mw.statusSbiHandle("导出电影数据失败!", mw.failureIco)
	}
}

func (mw *MyMainWindow) createMovieModel(searchInfo *SearchInfo) {
	var info, name string
	defer func() {
		mw.movieNameLb.SetCurrentIndex(-1)
	}()

	movieInfoMap := searchInfo.movieInfoMap

	mi := &MovieItems{movieItems: make([]MovieItem, len(movieInfoMap))}
	for i, movieInfo := range movieInfoMap {
		info = ""
		name = ""
		info = "第" + strconv.Itoa(i+1) + "部电影:\r\n"
		for key, value := range movieInfo {
			if key == "标　　题" && value != "-" {
				name = value
			}
			if key == "译　　名" && value != "-" {
				if name == "" {
					name = value
				}
			}
			if key == "片　　名" && value != "-" {
				if name == "" {
					name = value
				}
			}
			info += key + ":" + value + "'\r\n'"
		}
		if name == "" {
			name = "无名电影"
		}
		mi.movieItems[i] = MovieItem{name, info}
	}
	mm := &MovieModel{items: mi.movieItems}
	mw.movieNameLb.SetModel(mm)
	mw.movieModel = mm
	mw.statusTEHandle("搜索到" + strconv.Itoa(len(mm.items)) + "部电影")
}

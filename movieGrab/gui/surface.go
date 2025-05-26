package gui

import (
	"fmt"
	"log"
	"movieGrab/api"
	"movieGrab/controller"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var (
	wg       sync.WaitGroup
	dygHan   *controller.DygangHandler
	dytt8Han *controller.Dytt8Handler
	ys66Han  *controller.Ys66Handler
	allHan   *controller.AllSearchMovieHandler
)

type MyDialog struct {
	*walk.Dialog
	ui myDialogUI
}

type SearchInfo struct {
	MovieWeb          string
	MovieType         int
	MoviePage         string
	ResultTotalPages  int
	ResultTotalMovies int
	movieInfoMap      map[int]map[string]string
}

type MyMainWindow struct {
	*walk.MainWindow
	prevFilePath       string
	movieNameLb        *walk.ListBox
	movieInfoTE        *walk.TextEdit
	statusTE           *walk.TextEdit
	allMovieTE         *walk.TextEdit
	nowtimeTE          *walk.Menu
	downloadNameTE     *walk.TextEdit
	downloadUrlTE      *walk.TextEdit
	downloadSavePathTE *walk.TextEdit
	movieInfoLL        *walk.LinkLabel
	onlineWV           *walk.WebView
	dygMovie           *walk.RadioButton
	dyttMovie          *walk.RadioButton
	dy66Movie          *walk.RadioButton
	allMovie           *walk.RadioButton
	pagesLable         *walk.Label
	pagesNE            *walk.TextEdit
	search             *walk.PushButton
	searchAll          *walk.PushButton
	contirm            *walk.PushButton
	export             *walk.PushButton
	download           *walk.PushButton
	cbMovieType        *walk.ComboBox
	movieInfoDB        *walk.DataBinder
	filePathDB         *walk.DataBinder
	downloadUrlDB      *walk.DataBinder
	tabWidget          *walk.TabWidget
	filePathDlg        *walk.Dialog
	yesOrNoDlg         *walk.Dialog
	MessageDlg         *walk.Dialog
	multiThreadDLDlg   *walk.Dialog
	downloadUrlDlg     *walk.Dialog
	baseLE             *walk.LineEdit
	modelLE            *walk.LineEdit
	saveLE             *walk.LineEdit
	failPathLE         *walk.LineEdit
	statusSbi          *walk.StatusBarItem
	progressSbi        *walk.StatusBarItem
	nowtimeSbi         *walk.StatusBarItem
	searchInfo         *SearchInfo
	filePath           *FilePath
	movieModel         *MovieModel
	failureIco         *walk.Icon
	successIco         *walk.Icon
	fileIco            *walk.Icon
	magnetCB           *walk.ComboBox
	quarkCB            *walk.ComboBox
	onlineCB           *walk.ComboBox
	magentChB          *walk.CheckBox
	quarkChB           *walk.CheckBox
	onlineChB          *walk.CheckBox
	magnetUrl          []*api.UrlType
	quarkUrl           []*api.UrlType
	onlineUrl          []*api.UrlType
}

func InitSurface() {
	mw := &MyMainWindow{}
	mw.filePath = new(FilePath)
	mw.searchInfo = new(SearchInfo)

	imgPath := getUserImageFilePath()

	mw.failureIco, _ = walk.NewIconFromFile(imgPath + "stop.ico")
	mw.successIco, _ = walk.NewIconFromFile(imgPath + "check.ico")
	mw.fileIco, _ = walk.Resources.Icon(imgPath + "file.ico")
	titleIcon, _ := walk.NewBitmapFromFile(imgPath + "titleIcon.png")
	backgroundImg, _ := walk.NewBitmapFromFile(imgPath + "backgroundImg.jpg")
	fileImg, _ := walk.NewBitmapFromFile(imgPath + "file.png")
	settingImg, _ := walk.NewBitmapFromFile(imgPath + "setting.png")
	exitImg, _ := walk.NewBitmapFromFile(imgPath + "exit.png")
	fileFolderImg, _ := walk.NewBitmapFromFile(imgPath + "fileFolder.png")
	downloadImg, _ := walk.NewBitmapFromFile(imgPath + "download.png")
	linkedImg, _ := walk.NewBitmapFromFile(imgPath + "linked.png")
	unkownImg, _ := walk.NewBitmapFromFile(imgPath + "unkown.png")
	aboutImg, _ := walk.NewBitmapFromFile(imgPath + "about.png")
	movieImg, _ := walk.NewBitmapFromFile(imgPath + "mov.png")
	xlsxImg, _ := walk.NewBitmapFromFile(imgPath + "xlsx.png")
	radioBackgroundImg, _ := walk.NewBitmapFromFile(imgPath + "radioBackground.jpg")
	textedit1Img, _ := walk.NewBitmapFromFile(imgPath + "textedit1.png")
	textedit2Img, _ := walk.NewBitmapFromFile(imgPath + "textedit2.png")

	MainWindow{
		Icon: titleIcon,
		Background: BitmapBrush{
			backgroundImg,
		},
		//GradientBrush{
		//	Vertexes: []walk.GradientVertex{
		//	{X: 0, Y: 0, Color: walk.RGB(255, 255, 127)},
		//	{X: 1, Y: 0, Color: walk.RGB(127, 191, 255)},
		//	{X: 0.5, Y: 0.5, Color: walk.RGB(255, 255, 255)},
		//	{X: 1, Y: 1, Color: walk.RGB(127, 255, 127)},
		//	{X: 0, Y: 1, Color: walk.RGB(255, 127, 127)},
		//},
		//Triangles: []walk.GradientTriangle{
		//	{0, 1, 2},
		//	{1, 3, 2},
		//	{3, 4, 2},
		//	{4, 0, 2},
		//},
		//},
		Title:    "电影助手",
		Size:     Size{800, 600},
		MinSize:  Size{800, 600},
		MaxSize:  Size{800, 600},
		AssignTo: &mw.MainWindow,
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:   &mw.movieInfoDB,
			DataSource: mw.searchInfo,
			Name:       "searchInfo",
		},

		MenuItems: []MenuItem{
			Menu{
				Text:  "&文件",
				Image: fileImg,
				Items: []MenuItem{
					Action{
						Text:  "&设置",
						Image: settingImg,
						OnTriggered: func() {
							if cmd, err := mw.setResultSavePathDialog(mw, mw.filePath); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								fmt.Println(walk.DlgCmdOK)
							}
						},
					},
					Action{
						Text:  "退出",
						Image: exitImg,
						OnTriggered: func() {
							os.Exit(200)
						},
					},
				},
			},
			Menu{
				Text:  "&工具",
				Image: fileFolderImg,
				Items: []MenuItem{
					Action{
						Text:  "多线程下载",
						Image: downloadImg,
						OnTriggered: func() {
							if cmd, err := mw.multiThreadDownloadMovie(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								mw.messageTigs("下载成功！")
							}
						},
					},
					Action{
						Text:  "下载链接",
						Image: linkedImg,
						OnTriggered: func() {
							index := mw.movieNameLb.CurrentIndex()
							if -1 == index {
								walk.MsgBox(mw, "电影链接", "请在电影列表中先选择需要查看的电影", walk.MsgBoxIconWarning)
							} else {
								mw.getMovieLinkDialog(mw, mw.searchInfo.movieInfoMap[index]["电影链接"])
							}
						},
					},
					Action{
						Text:  "导出表格",
						Image: xlsxImg,
						OnTriggered: func() {
							mw.setYesOrNoDialog(mw, "导出信息设置")
						},
					},
					Action{
						Text:  "在线美剧",
						Image: movieImg,
						OnTriggered: func() {
							mw.onlineMovie(mw, "https://www.meijutt.tv/")
						},
					},
				},
			},
			Menu{
				Text:  "&帮助",
				Image: unkownImg,
				Items: []MenuItem{
					Action{
						Text:        "关于",
						Image:       aboutImg,
						OnTriggered: mw.aboutActionTriggered,
					},
				},
			},
		},
		Children: []Widget{

			Composite{
				Name:    "功能框",
				MinSize: Size{500, 35},
				MaxSize: Size{500, 35},
				Layout: Grid{
					Columns: 9,
					Spacing: 5,
				},
				Children: []Widget{
					Label{
						TextColor: 16777215,
						Text:      "电影资源站: ",
					},

					RadioButtonGroup{
						DataMember: "MovieWeb",
						Buttons: []RadioButton{
							RadioButton{
								Name: "dygRadio",
								Background: BitmapBrush{
									radioBackgroundImg,
								},
								Text:     "电影港",
								Value:    "dygang",
								AssignTo: &mw.dygMovie,
								OnClicked: func() {
									mw.search.SetEnabled(true)
									mw.searchAll.SetEnabled(false)
									mw.cbMovieType.SetEnabled(true)
									mw.getMovieKinds()
								},
							},
							RadioButton{
								Name: "dyttRadio",
								Background: BitmapBrush{
									radioBackgroundImg,
								},
								Text:     "电影天堂",
								Value:    "dytt8",
								AssignTo: &mw.dyttMovie,
								OnClicked: func() {
									mw.search.SetEnabled(true)
									mw.searchAll.SetEnabled(false)
									mw.cbMovieType.SetEnabled(true)
									mw.getMovieKinds()
								},
							},
							RadioButton{
								Name: "dy66Radio",
								Background: BitmapBrush{
									radioBackgroundImg,
								},
								Text:     "66影视",
								Value:    "66ys",
								AssignTo: &mw.dy66Movie,
								OnClicked: func() {
									mw.search.SetEnabled(true)
									mw.searchAll.SetEnabled(false)
									mw.cbMovieType.SetEnabled(true)
									mw.getMovieKinds()
								},
							},
							RadioButton{
								Name: "allRadio",
								Background: BitmapBrush{
									radioBackgroundImg,
								},
								Text:     "所有",
								Value:    "all",
								AssignTo: &mw.allMovie,
								OnClicked: func() {
									mw.search.SetEnabled(false)
									mw.searchAll.SetEnabled(true)
									mw.cbMovieType.SetEnabled(false)
									mw.getMovieKinds()
								},
							},
						},
					},
					Label{
						Text:      "电影类型：",
						TextColor: 16777215,
					},
					ComboBox{
						Value:         Bind("MovieType", SelRequired{}),
						BindingMember: "TypeId",
						DisplayMember: "TypeName",
						AssignTo:      &mw.cbMovieType,
						MinSize:       Size{90, 0},
						MaxSize:       Size{90, 0},
						//OnCurrentIndexChanged: ,
					},
					Label{
						Text: "",
					},
					PushButton{
						Text:     "分类搜索",
						AssignTo: &mw.search,
						MaxSize:  Size{100, 20},
						MinSize:  Size{100, 20},
						OnClicked: func() {
							mw.statusTEHandle("电影类型统计中,请稍后...")
							mw.getMoviePages(mw.searchInfo)
						},
					},
				},
			},

			Composite{
				Name:    "综合搜索",
				MaxSize: Size{500, 80},
				Layout: Grid{
					Columns: 4,
					Spacing: 5,
				},
				Children: []Widget{
					Label{
						Text:      "关键字搜索: ",
						TextColor: 16777215,
					},
					TextEdit{
						AssignTo: &mw.allMovieTE,
						MaxSize:  Size{475, 20},
						MinSize:  Size{475, 20},
					},
					Label{},
					PushButton{
						Text:     "搜索",
						AssignTo: &mw.searchAll,
						MaxSize:  Size{100, 20},
						MinSize:  Size{100, 20},
						OnClicked: func() {
							if mw.allMovieTE.Text() == "" {
								walk.MsgBox(mw, "消息提示", "输入信息不能为空,请输入需要查询的关键字", walk.MsgBoxIconInformation)
								return
							}
							mw.statusTEHandle("电影类型统计中,请稍后...")
							mw.statusSbiHandle("全网关键词搜索中...", mw.failureIco)
							mw.searchAllMovie(mw.searchInfo)
						},
					},
				},
			},

			Composite{
				Layout: Grid{Columns: 3, Spacing: 10},
				Children: []Widget{
					ListBox{
						MaxSize:  Size{200, 400},
						MinSize:  Size{200, 400},
						AssignTo: &mw.movieNameLb,
						//Model:                 mw.newEnvModel(),
						OnCurrentIndexChanged: mw.lbCurrentIndexChanged,
						OnItemActivated:       mw.lbItemActivated,
					},
					GroupBox{
						//Title: "电影信息：",
						Layout:  VBox{},
						MinSize: Size{450, 400},
						Children: []Widget{
							Label{
								Text:      "信息：",
								TextColor: 16777215,
							},
							TextEdit{
								AssignTo: &mw.movieInfoTE,
								ReadOnly: true,
								MaxSize:  Size{500, 200},
								MinSize:  Size{500, 200},
								VScroll:  true,
								//TextColor: 16777215,
								Background: BitmapBrush{
									textedit1Img,
								},
							},
							LinkLabel{
								AssignTo: &mw.movieInfoLL,
								Visible:  true,
							},
							Label{
								Text:      "状态：",
								TextColor: 16777215,
							},
							TextEdit{
								AssignTo:  &mw.statusTE,
								ReadOnly:  true,
								MaxSize:   Size{500, 200},
								MinSize:   Size{500, 200},
								VScroll:   true,
								TextColor: 16777215,
								Background: BitmapBrush{
									textedit2Img,
								},
							},
						},
					},
				},
			},

			Composite{
				Layout: Grid{Columns: 6, Spacing: 9},
				Children: []Widget{
					Label{
						AssignTo:  &mw.pagesLable,
						Text:      "总共0部电影,共0页,选择下载页:",
						TextColor: 16777215,
					},
					TextEdit{
						AssignTo: &mw.pagesNE,
						MinSize:  Size{70, 0},
						MaxSize:  Size{70, 0},
						Text:     Bind("MoviePage"),
					},
					Label{
						Text:      "页（例:1-9将下载1到9页或输入3则下载第3页）",
						TextColor: 16777215,
					},
					PushButton{
						Text:     "查看",
						AssignTo: &mw.contirm,
						OnClicked: func() {
							mw.statusSbiHandle("获取电影过程中,请稍后...", mw.failureIco)
							mw.getMovieInfo(mw.searchInfo)
							mw.createMovieModel(mw.searchInfo)
						},
					},
					PushButton{
						Text:     "导出",
						AssignTo: &mw.export,
						OnClicked: func() {
							mw.setYesOrNoDialog(mw, "导出信息设置")
						},
					},
					PushButton{
						Text:     "下载",
						AssignTo: &mw.download,
						OnClicked: func() {
							mw.lbItemActivated()
							mw.messageTigs("下载成功！")
						},
					},
				},
			},
		},

		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				AssignTo: &mw.statusSbi,
				Icon:     mw.failureIco,
				Text:     "目前状态:等待接收指令...",
				Width:    360,
			},
			StatusBarItem{
				AssignTo: &mw.progressSbi,
				Width:    220,
				Text:     "当前下载进度:0.00%",
			},
			StatusBarItem{
				AssignTo: &mw.nowtimeSbi,
				Width:    200,
				Text:     "",
			},
		},
	}.Create()

	go func() {
		for {
			time.Sleep(time.Second)
			mw.nowtimeSbi.SetText(mw.getNowTime())
		}
	}()

	mw.dygMovie.SetFocus()
	mw.cbMovieType.SetCurrentIndex(0)
	mw.pagesNE.SetText("1")
	mw.contirm.SetEnabled(false)
	mw.download.SetEnabled(false)
	mw.export.SetEnabled(false)
	//固定窗体大小
	win.SetWindowLong(mw.Handle(), win.GWL_STYLE, win.GetWindowLong(mw.Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
	mw.getPathInfo()
	mw.Run()
}

// 打开文件功能
func (mw *MyMainWindow) openActionTriggered() {
	if err := mw.openDir(); err != nil {
		fmt.Println(err)
	}
	mw.setPathInfo()
}

// 关于
func (mw *MyMainWindow) aboutActionTriggered() {
	walk.MsgBox(mw, "关于", "FFVII Difa Design", walk.MsgBoxIconInformation)
}

func (mm *MovieModel) ItemCount() int {
	return len(mm.items)
}

func (mm *MovieModel) Value(index int) interface{} {
	return mm.items[index].name
}

// 下拉框
func (mw *MyMainWindow) lbCurrentIndexChanged() {
	i := mw.movieNameLb.CurrentIndex()
	if i < 0 {
		mw.download.SetEnabled(false)
		mw.export.SetEnabled(false)
		return
	}
	mw.download.SetEnabled(true)
	mw.export.SetEnabled(true)
	item := &mw.movieModel.items[i]
	mw.movieInfoTE.SetText(item.value)
}

func (mw *MyMainWindow) lbItemActivated() {
	magnet := `magnet\:\?xt\=.*`
	moviename := `dn=(.*?)&`
	quark := `https\:\/\/pan\.quark\.cn\/.*`
	online := `https\:\/\/www.*\.html`
	magnetCount := 0
	quarkCount := 0
	onlineCount := 0
	i := mw.movieNameLb.CurrentIndex()
	if i < 0 {
		return
	}
	movie := mw.searchInfo.movieInfoMap[i]
	for key, movieInfo := range movie {
		if "电影链接" == key {
			movieUrls := strings.Split(movieInfo, `*:`)
			mw.magnetUrl = make([]*api.UrlType, 0)
			mw.quarkUrl = make([]*api.UrlType, 0)
			mw.onlineUrl = make([]*api.UrlType, 0)

			mre := regexp.MustCompile(magnet)
			mName := regexp.MustCompile(moviename)
			qre := regexp.MustCompile(quark)
			ore := regexp.MustCompile(online)
			for _, movieUrl := range movieUrls {
				if mUrl := mre.FindString(movieUrl); "" != mUrl {
					temp := new(api.UrlType)
					temp.UrlId = magnetCount
					if len(mName.FindStringSubmatch(mUrl)) >= 1 {
						temp.UrlName, _ = url.QueryUnescape(mName.FindStringSubmatch(mre.FindString(movieUrl))[1])
					} else if len(mName.FindStringSubmatch(mUrl)) > 0 {
						temp.UrlName, _ = url.QueryUnescape(mName.FindStringSubmatch(mre.FindString(movieUrl))[0])
					} else {
						temp.UrlName = "无链接名称"
					}

					temp.UrlLink = mre.FindString(movieUrl)
					mw.magnetUrl = append(mw.magnetUrl, temp)
					magnetCount++
				}
				if "" != qre.FindString(movieUrl) {
					temp := new(api.UrlType)
					temp.UrlId = quarkCount
					temp.UrlName = "跳转到网盘"
					temp.UrlLink = qre.FindString(movieUrl)
					mw.quarkUrl = append(mw.quarkUrl, temp)
					quarkCount++
				}
				if "" != ore.FindString(movieUrl) {
					temp := new(api.UrlType)
					temp.UrlId = onlineCount
					temp.UrlName = "跳转到在线观看"
					temp.UrlLink = ore.FindString(movieUrl)
					mw.onlineUrl = append(mw.onlineUrl, temp)
					onlineCount++
				}
			}
		}
	}
	mw.selectDownloadUrl(mw)
}

// 状态消息框
func (mw *MyMainWindow) statusTEHandle(newMsg string) {
	oldMsg := mw.statusTE.Text()
	t := time.Now().Unix()
	tStr := time.Unix(t, 0).Format("2006-01-02 15:04:05")
	resMsg := oldMsg + "\r\n" + tStr + ":" + newMsg
	mw.statusTE.SetText(resMsg)
}

// 状态条框
func (mw *MyMainWindow) statusSbiHandle(newMsg string, iconMsg *walk.Icon) {
	mw.statusSbi.SetText(newMsg)
	mw.statusSbi.SetIcon(iconMsg)
}

package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"movieGrab/api"
	"movieGrab/utils"
	"movieGrab/utils/multiThreadDownload"
	"strings"
	"time"
)

func (mw *MyMainWindow) messageTigs(msg string) {
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ni.Dispose()
	icon, _ := walk.Resources.Icon("gui/img/stop.ico")

	if err = ni.SetIcon(icon); err != nil {
		fmt.Println(err)
		return
	}

	if err = ni.SetToolTip("tools"); err != nil {
		fmt.Println(err)
		return
	}

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		if err := ni.ShowCustom("电影助手", msg, icon); err != nil {
			fmt.Println(err)
			return
		}
	})

	exitAction := walk.NewAction()
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	ni.ContextMenu().Actions().Add(exitAction)
	ni.SetVisible(true)

	ni.ShowInfo("电影助手", msg)
}

// 在线播放界面
func (mw *MyMainWindow) onlineMovie(owner walk.Form, onlineUrl string) (int, error) {
	return Dialog{
		Title:   "Walk WebView Example'",
		MinSize: Size{800, 600},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			WebView{
				AssignTo: &mw.onlineWV,
				Name:     "wv",
				URL:      onlineUrl,
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}
				return "stop", nil
			},
		},
	}.Run(owner)
}

// 多线程下载界面
func (mw *MyMainWindow) multiThreadDownloadMovie(owner walk.Form) (int, error) {
	var acceptPB, cancelPB, clipboardPB *walk.PushButton
	return Dialog{
		AssignTo:      &mw.multiThreadDLDlg,
		Title:         "多线程下载",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{350, 200},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "请输入下载文件名:\r\n",
					},
					TextEdit{
						Text:     "",
						AssignTo: &mw.downloadNameTE,
						MaxSize:  Size{170, 20},
						MinSize:  Size{170, 20},
						//ReadOnly: true,
					},
					Label{
						Text: "请输入下载地址:\r\n",
					},
					TextEdit{
						Text:     "",
						AssignTo: &mw.downloadUrlTE,
						MaxSize:  Size{170, 20},
						MinSize:  Size{170, 20},
						//ReadOnly: true,
					},
					Label{
						Text: "下载保存路径:\r\n",
					},
					TextEdit{
						Text:     mw.filePath.BasePath + "\r\n",
						AssignTo: &mw.downloadSavePathTE,
						MaxSize:  Size{170, 20},
						MinSize:  Size{170, 20},
					},
					PushButton{
						Text:    "浏览...",
						MaxSize: Size{70, 20},
						MinSize: Size{70, 20},
						OnClicked: func() {
							dlg := new(walk.FileDialog)
							dlg.InitialDirPath = "."
							dlg.Title = "选择保存到文件夹路径"

							if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
								return
							} else if !ok {
								return
							}
							mw.downloadSavePathTE.SetText(dlg.FilePath)
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &clipboardPB,
						Text:     "粘贴",
						OnClicked: func() {
							content, err := walk.Clipboard().Text()
							if err != nil {
								panic(err)
							}
							mw.downloadUrlTE.SetText(content)
						},
					},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确定",
						OnClicked: func() {
							t := time.Now()
							dtc := multiThreadDownload.DownloadThreadController{ThreadCount: 9, FileUrl: mw.downloadUrlTE.Text(), DownloadFolder: mw.downloadSavePathTE.Text(), DownloadFileName: mw.downloadNameTE.Text()}
							go dtc.Download(1024 * 1024 * 2)
							mw.statusTEHandle("开始执行多线程下载任务...")
							mw.statusSbiHandle(dtc.DownloadStatus(), mw.failureIco)
							go func() {
								for {
									time.Sleep(time.Second)
									if "当前下载进度:100.00%" == dtc.DownloadProcessStatus() {

										mw.statusSbiHandle(dtc.DownloadStatus(), mw.successIco)
										break
									} else if "" != dtc.DownloadProcessStatus() {
										//bar := utils.NewBar(0, len(pageUrlChan))
										//bar.Add(1)
										mw.statusSbiHandle(dtc.DownloadStatus(), mw.failureIco)
										mw.progressSbi.SetText(dtc.DownloadProcessStatus())
									}
								}
							}()
							mw.statusSbiHandle(dtc.DownloadStatus()+fmt.Sprintf("总共花费时间:%v", time.Since(t)), mw.successIco)
							mw.multiThreadDLDlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "取消",
						OnClicked: func() {
							mw.multiThreadDLDlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}

// 是否对话框
func (mw *MyMainWindow) getMovieLinkDialog(owner walk.Form, movieLink string) (int, error) {
	var acceptPB *walk.PushButton
	movieLink = utils.MovieLinkHandle(movieLink)
	return Dialog{
		AssignTo:      &mw.MessageDlg,
		Title:         "电影链接",
		DefaultButton: &acceptPB,
		MinSize:       Size{400, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					TextEdit{
						Text:    movieLink,
						VScroll: true,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确定",
						OnClicked: func() {
							mw.MessageDlg.Accept()
						},
					},
				},
			},
		},
	}.Run(owner)
}

// 是否对话框
func (mw *MyMainWindow) setYesOrNoDialog(owner walk.Form, title string) (int, error) {
	var acceptPB, cancelPB *walk.PushButton
	return Dialog{
		AssignTo:      &mw.yesOrNoDlg,
		Title:         title,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{300, 200},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "是否修改当前保存路径?\r\n",
					},
					Label{
						Text: "\r\n",
					},
					Label{
						Text: "基本路径：" + mw.filePath.BasePath + "\r\n",
					},
					Label{
						Text: "\r\n",
					},
					Label{
						Text: "模板路径:" + mw.filePath.ModelPath + "\r\n",
					},
					Label{
						Text: "\r\n",
					},
					Label{
						Text: "导出路径:" + mw.filePath.SavePath + "\r\n",
					},
					Label{
						Text: "\r\n",
					},
					Label{
						Text: "未成功文件路径:" + mw.filePath.FailExportPath + "\r\n",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "修改后导出",
						OnClicked: func() {
							if cmd, err := mw.setResultSavePathDialog(mw, mw.filePath); err != nil {
								log.Print(err)
								return
							} else if cmd == walk.DlgCmdOK {
								mw.exportMovieToExcel(mw.searchInfo)
							}
							mw.yesOrNoDlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "直接导出",
						OnClicked: func() {
							mw.exportMovieToExcel(mw.searchInfo)
							mw.yesOrNoDlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}

// 设置保存路径窗口
func (mw *MyMainWindow) setResultSavePathDialog(owner walk.Form, filePath *FilePath) (int, error) {
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &mw.filePathDlg,
		Title:         "设置导出文件保存路径",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &mw.filePathDB,
			Name:           "FilePathSet",
			DataSource:     filePath,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "基本路径:",
					},
					LineEdit{
						AssignTo: &mw.baseLE,
						Text:     Bind("BasePath"),
					},
					Label{
						Text: "模板路径:",
					},
					LineEdit{
						AssignTo: &mw.modelLE,
						Text:     Bind("ModelPath"),
					},
					Label{
						Text: "导出路径:",
					},
					LineEdit{
						AssignTo: &mw.saveLE,
						Text:     Bind("SavePath"),
					},
					Label{
						Text: "未成功文件路径:",
					},
					LineEdit{
						AssignTo: &mw.failPathLE,
						Text:     Bind("FailExportPath"),
					},
					PushButton{
						Text:      "浏览...",
						OnClicked: mw.openActionTriggered,
					},
					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确定",
						OnClicked: func() {
							if err := mw.filePathDB.Submit(); err != nil {
								log.Print(err)
								return
							}
							mw.filePathDlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { mw.filePathDlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

// 选择下载方式
func (mw *MyMainWindow) selectDownloadUrl(owner walk.Form) (int, error) {
	var urlType *api.UrlType
	var acceptPB, cancelPB *walk.PushButton
	return Dialog{
		AssignTo: &mw.downloadUrlDlg,
		Name:     "selectDownloadUrl",
		Title:    "选择下载链接",
		//DefaultButton: &acceptPB,
		//CancelButton:  &cancelPB,
		MinSize: Size{400, 200},
		DataBinder: DataBinder{
			AssignTo:       &mw.downloadUrlDB,
			Name:           "downloadUrlSet",
			DataSource:     urlType,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{
						Text: "种子下载(迅雷|电驴...):\r\n",
					},
					ComboBox{
						Value:         Bind("UrlLink", SelRequired{}),
						BindingMember: "UrlId",
						DisplayMember: "UrlName",
						AssignTo:      &mw.magnetCB,
						MinSize:       Size{150, 0},
						MaxSize:       Size{150, 0},
						Model:         mw.getMovieDownloadUrl("magnet"),
						Editable:      false,
					},
					CheckBox{
						AssignTo: &mw.magentChB,
						OnClicked: func() {
							mw.magnetCB.SetEnabled(true)
							mw.quarkCB.SetEnabled(false)
							mw.onlineCB.SetEnabled(false)
							mw.quarkChB.SetChecked(false)
							mw.onlineChB.SetChecked(false)
							mw.magentChB.SetEnabled(false)
							mw.quarkChB.SetEnabled(true)
							mw.onlineChB.SetEnabled(true)
						},
					},
					Label{
						Text: "网盘下载(夸克|迅雷云...):\r\n",
					},
					ComboBox{
						Value:         Bind("UrlLink", SelRequired{}),
						BindingMember: "UrlId",
						DisplayMember: "UrlName",
						AssignTo:      &mw.quarkCB,
						MinSize:       Size{150, 0},
						MaxSize:       Size{150, 0},
						Model:         mw.getMovieDownloadUrl("quark"),
						Editable:      false,
					},
					CheckBox{
						AssignTo: &mw.quarkChB,
						OnClicked: func() {
							mw.magnetCB.SetEnabled(false)
							mw.quarkCB.SetEnabled(true)
							mw.onlineCB.SetEnabled(false)
							mw.magentChB.SetChecked(false)
							mw.onlineChB.SetChecked(false)
							mw.magentChB.SetEnabled(true)
							mw.quarkChB.SetEnabled(false)
							mw.onlineChB.SetEnabled(true)
						},
					},
					Label{
						Text: "在线观看:\r\n",
					},
					ComboBox{
						Value:         Bind("UrlLink", SelRequired{}),
						BindingMember: "UrlId",
						DisplayMember: "UrlName",
						AssignTo:      &mw.onlineCB,
						MinSize:       Size{150, 0},
						MaxSize:       Size{150, 0},
						Model:         mw.getMovieDownloadUrl("online"),
						Editable:      false,
					},
					CheckBox{
						AssignTo: &mw.onlineChB,
						OnClicked: func() {
							mw.magnetCB.SetEnabled(false)
							mw.quarkCB.SetEnabled(false)
							mw.onlineCB.SetEnabled(true)
							mw.magentChB.SetChecked(false)
							mw.quarkChB.SetChecked(false)
							mw.magentChB.SetEnabled(true)
							mw.quarkChB.SetEnabled(true)
							mw.onlineChB.SetEnabled(false)
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						MaxSize:  Size{70, 0},
						AssignTo: &acceptPB,
						Text:     "确定",
						OnClicked: func() {
							//if err := mw.downloadUrlDB.Submit(); err != nil {
							//	log.Print(err)
							//	return
							//}
							magnetBool := mw.magnetCB.Enabled()
							quarkBool := mw.quarkCB.Enabled()
							onlineBool := mw.onlineCB.Enabled()
							if magnetBool {
								mindex := mw.magnetCB.CurrentIndex()
								mmodel := mw.magnetCB.Model()
								for i, m := range mmodel.([]*api.UrlType) {
									if i+1 == mindex {
										utils.OpenMagnet(m.UrlLink)
									}
								}
								mw.downloadUrlDlg.Accept()
							} else if quarkBool {
								mw.downloadUrlDlg.Accept()
							} else if onlineBool {
								oindex := mw.onlineCB.CurrentIndex()
								omodel := mw.onlineCB.Model()
								for i, m := range omodel.([]*api.UrlType) {
									if i == oindex {
										mw.onlineMovie(mw, m.UrlLink)
									}
								}
								mw.downloadUrlDlg.Accept()
							} else {
								walk.MsgBox(mw, "消息提示", "请选择一个链接对象", walk.DlgCmdClose)
							}
						},
					},
					PushButton{
						MaxSize:  Size{70, 0},
						AssignTo: &cancelPB,
						Text:     "取消",
						OnClicked: func() {
							mw.downloadUrlDlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}

func (mw *MyMainWindow) getMovieDownloadUrl(urlTpye string) []*api.UrlType {
	switch urlTpye {
	case "magnet":
		return mw.magnetUrl
	case "quark":
		return mw.quarkUrl
	case "online":
		return mw.onlineUrl
	default:
		return nil
	}
}

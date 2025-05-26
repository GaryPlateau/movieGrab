package gui

import (
	"github.com/lxn/walk"
	"movieGrab/config"
)

type FilePath struct {
	BasePath       string
	SavePath       string
	ModelPath      string
	FailExportPath string
}

func (mw *MyMainWindow) openDir() error {
	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.prevFilePath
	//dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
	dlg.InitialDirPath = "."
	dlg.Title = "选择保存到文件夹路径"

	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		return err
	} else if !ok {
		return nil
	}

	mw.prevFilePath = dlg.FilePath
	mw.baseLE.SetText(mw.prevFilePath)
	mw.modelLE.SetText(mw.prevFilePath + "\\model\\")
	mw.saveLE.SetText(mw.prevFilePath + "\\save\\")
	mw.failPathLE.SetText(mw.prevFilePath + "\\fail\\")

	mw.filePathDB.Submit()

	return nil
}

func (mw *MyMainWindow) setPathInfo() {
	conf := config.ConfigInfo{}
	filePath := make(map[string]string, 4)
	filePath["BasePath"] = mw.filePath.BasePath
	filePath["SavePath"] = mw.filePath.SavePath
	filePath["ModelPath"] = mw.filePath.ModelPath
	filePath["FailExportPath"] = mw.filePath.FailExportPath
	conf.CreatePath(filePath)
}

func (mw *MyMainWindow) getPathInfo() {
	conf := new(config.ConfigInfo)
	conf = conf.LoadConfig()
	mw.filePath.BasePath = conf.BasePath
	mw.filePath.SavePath = conf.SavePath
	mw.filePath.ModelPath = conf.ModelPath
	mw.filePath.FailExportPath = conf.FailExportPath
}

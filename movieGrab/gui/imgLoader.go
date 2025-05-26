package gui

import (
	"embed"
	"movieGrab/utils"
	"os"
)

//go:embed img/*
var fImg embed.FS

func getUserImageFilePath() (picPath string) {
	userPath, _ := os.UserHomeDir()
	picPath = userPath + "/AppData/Local/Temp/"
	return
}

func Imgload() {
	aboutImg, _ := fImg.ReadFile("img/about.png")
	checkIco, _ := fImg.ReadFile("img/check.ico")
	propertiesIco, _ := fImg.ReadFile("img/properties.ico")
	documentnewIco, _ := fImg.ReadFile("img/documentnew.ico")
	backgroundImg, _ := fImg.ReadFile("img/backgroundImg.jpg")
	settingPng, _ := fImg.ReadFile("img/setting.png")
	downloadPng, _ := fImg.ReadFile("img/download.png")
	movPng, _ := fImg.ReadFile("img/mov.png")
	movieJpg, _ := fImg.ReadFile("img/movie.jpg")
	openPng, _ := fImg.ReadFile("img/open.png")
	plusPng, _ := fImg.ReadFile("img/plus.png")
	exitPng, _ := fImg.ReadFile("img/exit.png")
	radioBackgroundJpg, _ := fImg.ReadFile("img/radioBackground.jpg")
	failureIco, _ := fImg.ReadFile("img/failure.ico")
	successIco, _ := fImg.ReadFile("img/success.ico")
	fileFolderPng, _ := fImg.ReadFile("img/fileFolder.png")
	filePng, _ := fImg.ReadFile("img/file.png")
	linkedPng, _ := fImg.ReadFile("img/linked.png")
	unkownPng, _ := fImg.ReadFile("img/unkown.png")
	uploadPng, _ := fImg.ReadFile("img/upload.png")
	xlsxPng, _ := fImg.ReadFile("img/xlsx.png")
	stopIco, _ := fImg.ReadFile("img/stop.ico")
	titleIcon, _ := fImg.ReadFile("img/titleIcon.png")
	textedit1Png, _ := fImg.ReadFile("img/textedit1.png")
	textedit2Png, _ := fImg.ReadFile("img/textedit2.png")

	imgPath := getUserImageFilePath()
	utils.CreatDir(imgPath)
	utils.WriteFile(imgPath+"titleIcon.png", titleIcon)
	utils.WriteFile(imgPath+"properties.ico", propertiesIco)
	utils.WriteFile(imgPath+"download.png", downloadPng)
	utils.WriteFile(imgPath+"mov.png", movPng)
	utils.WriteFile(imgPath+"movie.jpg", movieJpg)
	utils.WriteFile(imgPath+"open.png", openPng)
	utils.WriteFile(imgPath+"plus.png", plusPng)
	utils.WriteFile(imgPath+"about.png", aboutImg)
	utils.WriteFile(imgPath+"radioBackground.jpg", radioBackgroundJpg)
	utils.WriteFile(imgPath+"upload.png", uploadPng)
	utils.WriteFile(imgPath+"unkown.png", unkownPng)
	utils.WriteFile(imgPath+"xlsx.png", xlsxPng)
	utils.WriteFile(imgPath+"check.ico", checkIco)
	utils.WriteFile(imgPath+"backgroundImg.jpg", backgroundImg)
	utils.WriteFile(imgPath+"documentnew.ico", documentnewIco)
	utils.WriteFile(imgPath+"setting.png", settingPng)
	utils.WriteFile(imgPath+"exit.png", exitPng)
	utils.WriteFile(imgPath+"failure.ico", failureIco)
	utils.WriteFile(imgPath+"fileFolder.png", fileFolderPng)
	utils.WriteFile(imgPath+"file.png", filePng)
	utils.WriteFile(imgPath+"linked.png", linkedPng)
	utils.WriteFile(imgPath+"stop.ico", stopIco)
	utils.WriteFile(imgPath+"success.ico", successIco)
	utils.WriteFile(imgPath+"textedit1.png", textedit1Png)
	utils.WriteFile(imgPath+"textedit2.png", textedit2Png)
}

package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"movieGrab/utils"
	"os"
	"path/filepath"
)

//go:embed config.ini
var ConfigIni string

type ConfigInfo struct {
	BasePath       string `json:"BasePath"`
	ModelPath      string `json:"ModelPath"`
	SavePath       string `json:"SavePath"`
	FailExportPath string `json:"FailExportPath"`
}

func (this *ConfigInfo) CreatePath(filePathMap map[string]string) {
	exePath, _ := os.Executable()
	exePath = filepath.Dir(exePath)
	if len(filePathMap) == 0 {
		this.BasePath = exePath + "\\base\\"
		filePathMap["BasePath"] = this.BasePath
		this.ModelPath = exePath + "\\model\\"
		filePathMap["ModelPath"] = this.ModelPath
		this.SavePath = exePath + "\\save\\"
		filePathMap["SavePath"] = this.SavePath
		this.FailExportPath = exePath + "\\failExport\\"
		filePathMap["FailExportPath"] = this.FailExportPath
	} else {
		this.BasePath = filePathMap["BasePath"]
		this.ModelPath = filePathMap["ModelPath"]
		this.SavePath = filePathMap["SavePath"]
		this.FailExportPath = filePathMap["FailExportPath"]
	}

	jData, err := json.Marshal(filePathMap)
	if err != nil {
		fmt.Println("config序列化失败", err)
		return
	}
	ConfigIni = string(jData)
	//utils.WriteNewFile("config/config.ini", jData)
	utils.CreatDir(this.BasePath)
	utils.CreatDir(this.ModelPath)
	utils.CreatDir(this.SavePath)
	utils.CreatDir(this.FailExportPath)
	return
}

func (this *ConfigInfo) LoadConfig() *ConfigInfo {
	//configByte := utils.ReadFile("config/config.ini")
	if ConfigIni == "" {
		filePathMap := make(map[string]string, 4)
		this.CreatePath(filePathMap)
	}
	err := json.Unmarshal([]byte(ConfigIni), &this)
	if err != nil {
		fmt.Println("反序列化config失败:", err)
		return nil
	} else {
		utils.CreatDir(this.BasePath)
		utils.CreatDir(this.ModelPath)
		utils.CreatDir(this.SavePath)
		utils.CreatDir(this.FailExportPath)
	}
	return this
}

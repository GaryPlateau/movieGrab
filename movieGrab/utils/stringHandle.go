package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func MovieInfo66ysHander(movieInfo string) map[string]string {
	flag := checkCharacter(movieInfo)
	var movieMap map[string]string
	var movieInfoArr []string
	if flag {
		movieInfoArr = strings.FieldsFunc(movieInfo, splitChar)
	} else {
		movieInfoArr = strings.FieldsFunc(movieInfo, splitN)
	}

	infoLength := len(movieInfoArr)
	if 0 == infoLength {
		fmt.Println("获取电影信息失败")
		return nil
	}
	movieMap = make(map[string]string, infoLength)
	pattern := "([0-9A-Za-z\u4e00-\u9fa5]+[\u3000]*[0-9A-Za-z\u4e00-\u9fa5]+?[\u3000])" + `([\s\S]*` + "[\u3000]*" + `.+[\s\S])*`
	re := regexp.MustCompile(pattern)

	for _, val := range movieInfoArr {
		res := re.FindStringSubmatch(val)
		pattern1 := `[\s]+` + "[\u3000]+"
		re1 := regexp.MustCompile(pattern1)
		if len(res) != 0 {
			keyName := []rune(res[1])
			movieMap[string(keyName[:len(keyName)-1])] = re1.ReplaceAllString(res[2], "\r\n")
		}
	}

	return movieMap
}

func MovieInfoDytt8Hander(movieInfo string) map[string]string {
	var result map[string]string
	regtx := "◎"
	regTitle := "([0-9A-Za-z\u4e00-\u9fff]+[　]*[0-9A-Za-z\u4e00-\u9fff])[\\s\\r\\n　](.+)"
	regLink := "([电影链接]{4})(.+)"
	regSummary := "([\u4e00-\u9FFF][　]+介)"
	regContent := "([\u4e00-\u9FFF]{2}.+)"

	retx, _ := regexp.Compile(regtx)
	movieInfoArr := retx.Split(movieInfo, -1)
	result = make(map[string]string, len(movieInfoArr))

	reTitle, _ := regexp.Compile(regTitle)
	reLink, _ := regexp.Compile(regLink)
	reSum, _ := regexp.Compile(regSummary)
	reCont, _ := regexp.Compile(regContent)

	for i, info := range movieInfoArr {
		if i == 0 {
			continue
		}
		info = strings.Trim(info, "\r")
		info = strings.Trim(info, "　")
		infoT := reTitle.FindStringSubmatch(info)
		if i == len(movieInfoArr)-1 {
			infoSum := reSum.FindString(info)
			infoCont := reCont.FindString(info)
			result[infoSum] = infoCont
			infoT = reLink.FindStringSubmatch(info)
		}
		if len(infoT) > 1 {
			result[infoT[1]] = strings.Trim(infoT[2], " ")
		}
	}
	return result
}

// 整理电影信息到map
func MovieInfoHander(movieInfo string) map[string]string {
	flag := checkCharacter(movieInfo)
	var movieMap map[string]string
	var movieInfoArr []string
	if flag {
		movieInfoArr = strings.FieldsFunc(movieInfo, splitChar)
	} else {
		movieInfoArr = strings.FieldsFunc(movieInfo, splitN)
	}

	infoLength := len(movieInfoArr)
	if 0 == infoLength {
		fmt.Println("获取电影信息失败")
		return nil
	}
	movieMap = make(map[string]string, infoLength)
	pattern := "([0-9A-Za-z\u4e00-\u9fa5]+[\u3000]*[\u4e00-\u9fa5]+)" + `[\s]*` + "[\u3000]*" + "(.+)"
	re := regexp.MustCompile(pattern)
	for i := 0; i < infoLength; i++ {
		result := re.FindStringSubmatch(movieInfoArr[i])
		if result == nil {
			continue
		}
		if infoLength-2 == i {
			_, ok := movieMap["简　　介"]
			if !ok {
				movieMap["简　　介"] = movieInfoArr[i]
			}
		}
		if flag {
			movieMap[result[1]] = result[2]
		} else {
			movieMap[result[1]] = result[2][2:]
		}
		//movieMap[result[1]] = result[2]
	}
	//movieMap["简介"] = movieInfoArr[infoLength-2]
	return movieMap
}

func checkCharacter(cont string) bool {
	pattern := "◎"
	re := regexp.MustCompile(pattern)
	flag := re.MatchString(cont)
	return flag
}

func splitChar(c rune) bool {
	if c == '◎' {
		return true
	} else {
		return false
	}
}

func splitN(c rune) bool {
	if c == '\n' || c == '\r' {
		return true
	} else {
		return false
	}
}

func MovieLinkHandle(movielink string) string {
	reg := `\*\:`
	re := regexp.MustCompile(reg)
	return re.ReplaceAllString(movielink, "\r\n")
}

func clearMessyCode() {
	//var errorCode = []string{"聽","鈥","︹","€Α","禭"}
}

package gui

import "time"

func (mw *MyMainWindow) getNowTime() string {
	t := time.Now().Unix()
	tStr := time.Unix(t, 0).Format("2006-01-02 15:04:05")
	return tStr
}

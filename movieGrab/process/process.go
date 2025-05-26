package process

import (
	"movieGrab/controller"
	"sync"
)

var (
	wg            sync.WaitGroup
	htmlHander    controller.HtmlHandler
	pageTotalChan chan string
	pageUrlChan   chan string
	moiveInfoChan chan string
)

// 启动服务
func StartProcess(exitChan chan bool) {
	
}

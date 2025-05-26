package multiThreadDownload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"movieGrab/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DownloadTask struct {
	customFunc func(params interface{}) // 执行方法
	paramsInfo interface{}              // 执行方法参数
}

type DownloadThreadController struct {
	TaskQueue              chan DownloadTask       // 用于接收下载任务
	TaskCount              chan int                // 用于记载当前任务数量
	Exit                   chan int                // 用于记载当前任务数量
	ThreadCount            int                     // 最大协程数
	WaitGroup              sync.WaitGroup          // 等待协程完成
	RangeStrs              map[int]string          // 所有需要下载的文件名
	FileUrl                string                  // 下载链接
	DownloadResultInfoChan chan DownloadFileParams // 下载任务响应通道
	DownloadFolder         string                  // 下载文件保存文件夹
	DownloadFileName       string                  // 下载文件保存文件名
	Filenames              []string                // 子文件名，有序
	ProcessStatus          string                  // 下载进度
	DownloadMsg            string                  // 下载状态
}

type DownloadFileParams struct {
	UrlStr       string
	RangeStr     string
	RangeIndex   int
	TempFilename string
	Successed    bool
}

func (this *DownloadThreadController) Put(task DownloadTask) {
	// 用于开启单个协程任务，下载文件的部分内容
	defer func() {
		err := recover() //内置函数，可以捕捉到函数异常
		if err != nil {
			fmt.Println("Channel closed", err)
		}
	}()
	this.WaitGroup.Add(1)  // 每插入一个任务，就需要计数
	this.TaskCount <- 1    // 含缓冲区的通道，用于控制下载器的协程最大数量
	this.TaskQueue <- task // 插入下载任务
	//go task.customFunc(task.paramsInfo)
}

func (this *DownloadThreadController) DownloadProcessStatus() string {
	return "当前下载进度:" + this.ProcessStatus
}

func (this *DownloadThreadController) DownloadStatus() string {
	return this.DownloadMsg
}

func (this *DownloadThreadController) DownloadFile(paramsInfo interface{}) {
	// 下载任务，接收对应的参数，负责从网页中下载对应部分的文件资源
	defer func() {
		this.WaitGroup.Done() // 下载任务完成，协程结束
	}()
	switch paramsInfo.(type) {
	case DownloadFileParams:
		params := paramsInfo.(DownloadFileParams)
		params.Successed = false
		defer func() {
			err := recover() //内置函数，可以捕捉到函数异常
			if err != nil {
				// 如果任意环节出错，表明下载流程未成功完成，标记下载失败
				this.DownloadMsg = fmt.Sprintf("下载文件失败:%v", err)
				params.Successed = false
			}
		}()
		//fmt.Println("Start to down load " + params.UrlStr + ", Content-type: " + params.RangeStr + " , save to file: " + params.TempFilename)
		urlStr := params.UrlStr
		rangeStr := params.RangeStr
		tempFilename := params.TempFilename
		os.Remove(tempFilename) // 删除已有的文件, 避免下载的数据被污染
		// 发起文件下载请求
		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Header.Add("Range", rangeStr)      // 测试下载部分内容
		res, err := http.DefaultClient.Do(req) // 发出下载请求，等待回应
		if err != nil {
			this.DownloadMsg = fmt.Sprintf("链接目标地址失败:%v", urlStr)
			params.Successed = false // 无法连接, 标记下载失败
		} else if res.StatusCode != 206 {
			params.Successed = false
		} else { // 能正常发起请求
			// 打开文件，写入文件
			fileObj, err := os.OpenFile(tempFilename, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				this.DownloadMsg = fmt.Sprintf("打开临时文件失败 " + tempFilename)
				params.Successed = false // 无法打开文件, 标记下载失败
			} else {
				defer fileObj.Close()                 // 关闭文件流
				body, err := ioutil.ReadAll(res.Body) // 读取响应体的所有内容
				if err != nil {
					this.DownloadMsg = fmt.Sprintf("目标链接相应失败:%v", err)
					params.Successed = false
				} else {
					defer res.Body.Close()  // 关闭连接流
					fileObj.Write(body)     // 写入字节数据到文件
					params.Successed = true // 成功执行到最后一步，则表示下载成功
				}
			}
		}
		this.DownloadResultInfoChan <- params // 将下载结果传入
	}
}

func (this *DownloadThreadController) Run() {
	// 只需要将待下载的请求发送一次即可，成功了会直接剔除，不成功则由接收方重试
	for rangeIndex, rangeStr := range this.RangeStrs {
		params := DownloadFileParams{UrlStr: this.FileUrl, RangeStr: rangeStr, TempFilename: this.DownloadFolder + "/" + rangeStr, RangeIndex: rangeIndex, Successed: true} // 下载参数初始化
		task := DownloadTask{this.DownloadFile, params}
		this.Put(task) // 若通道满了会阻塞，等待空闲时再下载
	}
}

func (this DownloadThreadController) GetSuffix(contentType string) string {
	suffix := ""
	contentTypes := map[string]string{
		"image/gif":                    "gif",
		"image/jpeg":                   "jpg",
		"application/x-img":            "img",
		"image/png":                    "png",
		"application/json":             "json",
		"application/pdf":              "pdf",
		"application/msword":           "word",
		"application/octet-stream":     "rar",
		"application/x-zip-compressed": "zip",
		"application/x-msdownload":     "exe",
		"video/mpeg4":                  "mp4",
		"video/avi":                    "avi",
		"audio/mp3":                    "mp3",
		"text/css":                     "css",
		"application/x-javascript":     "js",
		"application/vnd.android.package-archive": "apk",
	}
	for key, value := range contentTypes {
		if strings.Contains(contentType, key) {
			suffix = value
			break
		}
	}
	return suffix
}

func (this *DownloadThreadController) ResultProccess(trunkSize int) string {
	// 负责处理各个协程下载资源的结果， 若成功则从下载列表中剔除，否则重新将该任务Put到任务列表中；超过5秒便会停止
	MAX_RETRY_TIME := 100
	nowRetryTime := 0
	result_msg := ""
	for {
		select {
		case resultInfo := <-this.DownloadResultInfoChan:
			<-this.TaskCount          // 取出一个计数器，表示一个协程已经完成
			if resultInfo.Successed { // 成功下载该文件，清除文件名列表中的信息
				delete(this.RangeStrs, resultInfo.RangeIndex) // 删除任务队列中的该任务（rangeStr队列）
				this.ProcessStatus = strconv.FormatFloat((1.0-float64(len(this.RangeStrs))/float64(trunkSize))*100, 'f', 2, 64) + "%"
				//fmt.Println("Download progress -> " + strconv.FormatFloat((1.0-float64(len(this.RangeStrs))/float64(trunkSize))*100, 'f', 2, 64) + "%")
				if len(this.RangeStrs) == 0 {
					result_msg = "SUCCESSED"
					break
				}
			} else {
				nowRetryTime += 1
				if nowRetryTime > MAX_RETRY_TIME { // 超过最大的重试次数退出下载
					result_msg = "MAX_RETRY"
					break
				}
				task := DownloadTask{customFunc: this.DownloadFile, paramsInfo: resultInfo} // 重新加载该任务
				go this.Put(task)
			}
		case task := <-this.TaskQueue:
			function := task.customFunc
			go function(task.paramsInfo)
		case <-time.After(5 * time.Second):
			result_msg = "TIMEOUT"
			break
		}

		if result_msg == "MAX_RETRY" {
			this.DownloadMsg = fmt.Sprintf("The network is unstable, exceeding the maximum number of redownloads.")
			break
		} else if result_msg == "SUCCESSED" {
			this.DownloadMsg = fmt.Sprintf("下载文件成功!")
			break
		} else if result_msg == "TIMEOUT" {
			this.DownloadMsg = fmt.Sprintf("下载超市!")
			break
		}
	}

	close(this.TaskCount)
	close(this.TaskQueue)
	close(this.DownloadResultInfoChan)
	return result_msg
}

func (this *DownloadThreadController) testDownload(urlStr string, perThreadSize int) (int, map[int]string, []string, string, error) {
	// 尝试连接目标资源，目标资源是否可以使用多线程下载
	length := 0
	rangeMaps := make(map[int]string)
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Range", "bytes=0-1") // 测试下载部分内容
	res, err := http.DefaultClient.Do(req)
	contentType := ""
	rangeIndex := 1
	filenames := []string{}
	if err != nil {
		rangeMaps[rangeIndex] = urlStr
		return length, rangeMaps, filenames, contentType, errors.New("Failed to connet " + urlStr)
	}
	if res.StatusCode != 206 {
		rangeMaps[rangeIndex] = urlStr
		return length, rangeMaps, filenames, contentType, errors.New("Http status is not equal to 206!")
	}
	// 206表示响应成功，仅仅返回部分内容
	contentLength := res.Header.Get("Content-Range")
	contentType = res.Header.Get("Content-Type")
	total_length, err := strconv.Atoi(strings.Split(contentLength, "/")[1])
	if err != nil {
		return length, rangeMaps, filenames, contentType, errors.New("Can't calculate the content-length form server " + urlStr)
	}
	now_length := 0 // 记录byte偏移量
	for {
		if now_length >= total_length {
			break
		}
		var tempRangeStr string // 记录临时文件名
		if now_length+perThreadSize >= total_length {
			tempRangeStr = "bytes=" + strconv.Itoa(now_length) + "-" + strconv.Itoa(total_length-1)
			now_length = total_length
		} else {
			tempRangeStr = "bytes=" + strconv.Itoa(now_length) + "-" + strconv.Itoa(now_length+perThreadSize-1)
			now_length = now_length + perThreadSize
		}
		rangeMaps[rangeIndex] = tempRangeStr
		filenames = append(filenames, tempRangeStr)
		rangeIndex = rangeIndex + 1
	}
	return total_length, rangeMaps, filenames, contentType, nil
}

func (this *DownloadThreadController) Download(oneThreadDownloadSize int) bool {
	fsc := FileSuffixCheck{}
	this.DownloadMsg = fmt.Sprintf("尝试链接目标下载文件...")
	length, rangeMaps, tempFilenames, contentType, err := this.testDownload(this.FileUrl, oneThreadDownloadSize)
	this.DownloadMsg = fmt.Sprintf("下载文件总大小:" + strconv.FormatFloat(float64(length)/(1024.0*1024.0), 'f', 2, 64) + "M")
	if err != nil {
		this.DownloadMsg = fmt.Sprintf("此文件不支持多线程下载")
		return false
	}
	this.DownloadMsg = fmt.Sprintf("配置文件成功，开始下载目标文件...")
	this.InitConfig() // 初始化通道、分片等配置
	//oneThreadDownloadSize := 1024 * 1024 * 2 // 1024字节 = 1024bite = 1kb -> 2M
	oneThreadDownloadSize = 1024 * 1024 * 4 // 1024字节 = 1024bite = 1kb -> 4M
	filenames := []string{}
	for _, value := range tempFilenames {
		filenames = append(filenames, this.DownloadFolder+"/"+value)
	}
	fileSuffix := this.GetSuffix(contentType)
	filename := this.DownloadFileName // 获取文件下载名
	this.Filenames = filenames        //下载文件的切片列表
	this.RangeStrs = rangeMaps        // 下载文件的Range范围
	go this.Run()                     // 开始下载文件
	proccessResult := this.ResultProccess(len(rangeMaps))
	downloadResult := false // 定义下载结果标记
	if proccessResult == "SUCCESSED" {
		absoluteFilename := this.DownloadFolder + "/" + filename + "." + fileSuffix
		downloadResult = this.CombineFiles(filename + "." + fileSuffix)
		if downloadResult {
			newSuffix := fsc.GetFileType(fsc.GetBytesFile(absoluteFilename, 10))
			err = os.Rename(absoluteFilename, this.DownloadFolder+"/"+filename+"."+newSuffix)
			if err != nil {
				downloadResult = false
				this.DownloadMsg = fmt.Sprintf("文件合成成功, 重命名文件失败,默认文件名:%v", absoluteFilename)
			} else {
				this.DownloadMsg = fmt.Sprintf("文件合成成功, 重命名文件成功, 新文件名:%v", this.DownloadFolder+"/"+filename+"."+newSuffix)
			}
		} else {
			this.DownloadMsg = fmt.Sprintf("下载文件失败")
		}
	} else {
		this.DownloadMsg = fmt.Sprintf("下载文件失败,错误原因:%v", proccessResult)
		downloadResult = false
	}
	return downloadResult
}

func (this *DownloadThreadController) CombineFiles(filename string) bool {
	os.Remove(this.DownloadFolder + "/" + filename)
	goalFile, err := os.OpenFile(this.DownloadFolder+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		this.DownloadMsg = fmt.Sprintf("打开文件失败: %v", err)
		return false
	}

	// 正确的话应按照初始计算的文件名顺序合并，并且无缺失
	for _, value := range this.Filenames {
		retryTime := 3
		tempFileBytes := []byte{}
		for retryTime > 0 {
			tempFileBytes = utils.ReadFile(value)
			time.Sleep(100) // 休眠100毫秒，看看是不是文件加载错误
			if tempFileBytes != nil {
				break
			}
			retryTime = retryTime - 1
		}
		goalFile.Write(tempFileBytes)
		os.Remove(value)
	}
	goalFile.Close()
	return true
}

func (this *DownloadThreadController) InitConfig() {
	taskQueue := make(chan DownloadTask, this.ThreadCount)
	taskCount := make(chan int, this.ThreadCount+1)
	exit := make(chan int)
	downloadResultInfoChan := make(chan DownloadFileParams)
	this.TaskQueue = taskQueue
	this.TaskCount = taskCount
	this.Exit = exit
	this.DownloadResultInfoChan = downloadResultInfoChan
	this.WaitGroup = sync.WaitGroup{}
	this.RangeStrs = make(map[int]string)
	utils.SafeMkdir(this.DownloadFolder)
}

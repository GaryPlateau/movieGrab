package utils

import (
	"bytes"
	"fmt"
	"sync"
)

func CreateNumbers(bitNum int, wg *sync.WaitGroup) {
	defer wg.Done()
	var buffer bytes.Buffer
	var mutex sync.Mutex
	for j := bitNum * 10000000; j < (bitNum+1)*10000000; j++ {
		tmp := fmt.Sprintf("%08d\n", j)
		buffer.Write([]byte(tmp))
	}
	contents := buffer.Bytes()
	mutex.Lock()
	WriteFile("D:/test.txt", contents)
	mutex.Unlock()
	fmt.Printf("完成%d数据写入\r", bitNum)
}

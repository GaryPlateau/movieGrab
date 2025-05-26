package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// O_RDONLY int = syscall.O_RDONLY // 只读模式打开文件
// O_WRONLY int = syscall.O_WRONLY // 只写模式打开文件
// O_RDWR   int = syscall.O_RDWR   // 读写模式打开文件
// O_APPEND int = syscall.O_APPEND // 写操作时将数据附加到文件尾部
// O_CREATE int = syscall.O_CREAT  // 如果不存在将创建一个新文件
// O_EXCL   int = syscall.O_EXCL   // 和O_CREATE配合使用，文件必须不存在
// O_SYNC   int = syscall.O_SYNC   // 打开文件用于同步I/O
// O_TRUNC  int = syscall.O_TRUNC  // 如果可能，打开时清空文件

// 读取文件 (大文件使用)
func ReadBigFile(filePath string) []byte {
	var buffer bytes.Buffer
	flag, _ := pathExists(filePath)
	if !flag {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return nil
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		lineByte, err := reader.ReadBytes('\n')
		buffer.Write(lineByte)
		if err == io.EOF {
			break
		}
	}
	result := buffer.Bytes()
	return result
}

// 一次性读取文件（小文件使用）
func ReadFile(filePath string) []byte {
	flag, _ := pathExists(filePath)
	if !flag {
		return nil
	}

	tempFile, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open file " + filePath)
		return nil
	}

	content, err := ioutil.ReadAll(tempFile)
	if err != nil {
		fmt.Println("读取文件内容失败"+filePath, err)
		return nil
	}
	tempFile.Close()
	return content
}

// 写byte内容到filePath中
func WriteFile(filePath string, content []byte) {

	fileBytes, _ := os.ReadFile(filePath)
	if fileBytes != nil {
		//fmt.Println("文件已经存在")
		return
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer file.Close()

	//使用写入缓存的方式
	writen := bufio.NewWriter(file)
	count, err := writen.Write(content)

	//需要使用Flush()将写入到writer缓存的数据真正写入到.txt文件中
	writen.Flush()
	if err != nil {
		fmt.Println("写入内容失败", err)
		fmt.Printf("写入%v字节数据", count)
		return
	} else {
		//fmt.Printf("写入%v字节数据", count)
		return
	}
}

func WriteNewFile(filePath string, content []byte) {

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer file.Close()

	//使用写入缓存的方式
	writen := bufio.NewWriter(file)
	count, err := writen.Write(content)

	//需要使用Flush()将写入到writer缓存的数据真正写入到.txt文件中
	writen.Flush()
	if err != nil {
		fmt.Println("写入内容失败", err)
		fmt.Printf("写入%v字节数据", count)
		return
	} else {
		//fmt.Printf("写入%v字节数据", count)
		return
	}
}

// 复制一个文件内容到另一个文件
func CopyFile(srcFile string, dscFile string) {
	src, _ := os.OpenFile(srcFile, os.O_RDONLY, 0777)
	defer src.Close()
	dsc, _ := os.OpenFile(srcFile, os.O_WRONLY, 0777)
	defer dsc.Close()
	writen, err := io.Copy(dsc, src)
	if err != nil {
		fmt.Println("复制文件出错", err)
		return
	} else {
		fmt.Printf("成功复制%v字节文件内容", writen)
		return
	}
}

// 检查文件是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 创建文件夹
func CreatDir(path string) {
	exist, err := pathExists(path)
	if err != nil {
		fmt.Printf("获取文件夹异常 -> %v\n", err)
		return
	}
	if exist {
		//fmt.Println("文件夹已存在！")
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("创建%v目录异常 -> %v\n", path, err)
		} else {
			//fmt.Println("创建成功!")
		}
	}
}

func SafeMkdir(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}

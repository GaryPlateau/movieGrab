package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func OpenMagnet(magentUrl string) {
	magentUrl = strings.Trim(magentUrl, " ")
	cmd := exec.Command("cmd", "/c", "start", magentUrl)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

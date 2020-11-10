package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/kardianos/osext"
	"github.com/webview/webview"
)

func main() {
	//启动服务
	cmd := exec.Command("/usr/bin/sh", "/home/xuan/ccal/web/v0.0.6/listen.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("农历择日")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://127.0.0.1:9090")
	w.Run()
}

func filePath() {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(folderPath)
}

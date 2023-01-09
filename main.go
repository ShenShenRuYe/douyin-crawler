package main

import (
	"bufio"
	"douyin/crawler"
	"douyin/utils"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	file := "log" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetPrefix("[抖音爬取工具]")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//create the folder to save video
	s := utils.SaveFolder{}
	s.FolderType = "video"
	s.Create()
}

func main() {
	//get the video
	var requestVideo = crawler.RequestVideo{}
	args := os.Args
	if len(args) != 1 {
		requestVideo.ShareText = args[1]
	} else {
		for {
			fmt.Println("请输入分享链接")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			userInput := input.Text()
			if userInput == "-1" {
				break
			}
			requestVideo.ShareText = input.Text()
			requestVideo.Save()
		}

	}

}

/*
Package utils
@Time : 2023/1/9 19:58
@Author : 董胜烨
@File : util
@Software: GoLand
@note:
*/
package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type SaveFolder struct {
	filePath   string
	FolderType string
}

func (receiver SaveFolder) FilePath() string {
	return receiver.filePath
}

func (receiver *SaveFolder) Create() {
	var message string
	// judge if is to create date-folder then generate the path of the folder
	if receiver.FolderType == "dateVideo" {
		date := fmt.Sprintf("%s-%s-%s", strconv.Itoa(time.Now().Year()), time.Now().Format("01"), time.Now().Format("02"))
		receiver.filePath = fmt.Sprintf("%s/%s", "videoFiles", date)
		message = date + "当日存储视频文件夹创建成功"
	} else if receiver.FolderType == "video" {
		receiver.filePath = "videoFiles"
		message = "视频文件夹创建成功"
	}
	// check is the folder exists
	exists := receiver.pathExists()
	if exists {
		return
	}
	// if not exists,then create
	err := os.Mkdir(receiver.filePath, os.ModePerm)
	if err != nil {
		return
	}
	log.Println(message)
}

func (receiver SaveFolder) pathExists() bool {
	_, err := os.Stat(receiver.filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

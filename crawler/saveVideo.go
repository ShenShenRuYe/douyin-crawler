/*
Package crawler
@Time : 2023/1/9 10:40
@Author : 董胜烨
@File : saveVideo
@Software: GoLand
@note:
*/
package crawler

import (
	"douyin/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type RequestVideo struct {
	ShareText string //分享语句

	shortURL string //短连接
	videoURL string //视频连接

	uid        string //视频唯一id
	videoTitle string //视频标题

	filePath string //保存路径
}

func (v *RequestVideo) getShortText() {
	re, _ := regexp.Compile(`https?://\S+`)
	matches := re.FindAllString(v.ShareText, -1)
	v.shortURL = matches[0]
}

func (v *RequestVideo) getUid() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", v.shortURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 获取视频的 id，进行短连接的重定向
	midURL, err := url.Parse(resp.Request.URL.String())
	if err != nil {
		return err
	}
	re, err := regexp.Compile("video/(\\d*)")

	if err != nil {
		return err
	}
	matches := re.FindStringSubmatch(midURL.Path)
	if len(matches) < 2 {
		return fmt.Errorf("Failed to parse unique id from URL: %s", midURL.Path)
	}
	uid := matches[1]
	v.uid = uid
	return nil
}

// getVideoURL 获取视频无水印的 url
func (v *RequestVideo) getVideoURL() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.iesdouyin.com/aweme/v1/web/aweme/detail/?aweme_id=%s&aid=1128&version_name=23.5.0&device_platform=android&os_version=2333", v.uid), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	title, src, err := parseVideoURLResponse(resp.Body)
	if err != nil {
		return err
	}
	v.videoTitle = title
	v.videoURL = src
	return nil
}

func parseVideoURLResponse(body io.ReadCloser) (string, string, error) {
	// 定义结构体用于解析响应 JSON
	type VideoURLResponse struct {
		AwemeDetail struct {
			Desc  string `json:"desc"`
			Video struct {
				PlayAddr struct {
					URLList []string `json:"url_list"`
				} `json:"play_addr"`
			} `json:"video"`
		} `json:"aweme_detail"`
	}
	// 解析响应
	var resp VideoURLResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return "", "", err
	}

	title := resp.AwemeDetail.Desc
	src := resp.AwemeDetail.Video.PlayAddr.URLList[0]
	return title, src, nil
}

// getVideo 获取视频
func (v *RequestVideo) getVideo() error {
	resp, err := http.Get(v.videoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	s := utils.SaveFolder{}
	s.FolderType = "dateVideo"
	s.Create()
	filePath := s.FilePath()
	out, err := os.Create(fmt.Sprintf("%s/%s.mp4", filePath, v.videoTitle))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func (v *RequestVideo) Save() {
	v.getShortText()
	err := v.getUid()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = v.getVideoURL()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = v.getVideo()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("【保存完成】%s！", v.videoTitle)
}

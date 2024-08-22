package fileup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
)

var (
	client = http.DefaultClient
)

type FileInfo struct {
	Filename      string
	ContentType   string
	ContentLength int64
	Link          string
}

func UploadByURL(ctx context.Context, fileURL, fileType string, params map[string]string) (info *FileInfo, err error) {
	u, err := url.ParseRequestURI(fileURL)
	if err != nil {
		return nil, errors.New("bad url")
	}

	downResp, err := client.Get(fileURL)
	if err != nil {
		return
	}
	defer downResp.Body.Close()

	log.Println("outer content-type:", downResp.ContentLength)
	log.Println("content-type:", downResp.Header.Get("content-type"))
	log.Println("content-length:", downResp.Header.Get("content-length"))
	log.Println("content-disposition:", downResp.Header.Get("content-disposition"))

	fileName := filepath.Base(u.Path)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	log.Println("form-data content-type", writer.FormDataContentType())
	// return

	part, err := writer.CreateFormFile("media", fileName)
	if err != nil {
		return
	}

	for k, v := range params {
		_ = writer.WriteField(k, v)
	}
	_, err = io.Copy(part, downResp.Body)
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	uploadURL := fmt.Sprintf("https://file.io/?title=%s", fileName)
	// 上传
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", "hikari-uploader")
	for k, v := range req.Header {
		log.Printf("%s  ->  %v\n", k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("upload resp fail: %v\n", err)
	}
	defer resp.Body.Close()

	log.Println("upload resp status", resp.Status)

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read upload resp fail: %v\n", err)
	}

	/*
			{
		    "success": true,
		    "status": 200,
		    "id": "ba4072d0-602f-11ef-8f04-dd5fa071332c",
		    "key": "3qPPi6iYNzpo",
		    "path": "/",
		    "nodeType": "file",
		    "name": "rpSjWSDDlCWaefl3.mp4",
		    "title": "rpSjWSDDlCWaefl3.mp4",
		    "description": null,
		    "size": 248745,
		    "link": "https://file.io/3qPPi6iYNzpo",
		    "private": false,
		    "expires": "2024-09-05T02:39:17.618Z",
		    "downloads": 0,
		    "maxDownloads": 1,
		    "autoDelete": true,
		    "planId": 0,
		    "screeningStatus": "pending",
		    "mimeType": "application/octet-stream",
		    "created": "2024-08-22T02:39:17.618Z",
		    "modified": "2024-08-22T02:39:17.618Z"
		}
	*/

	log.Println("upload resp:", string(bodyData))
	for k, v := range resp.Header {
		fmt.Printf("%s = %v\n", k, v)
	}

	var result struct {
		Link string `json:"link"`
		Size int    `json:"size"`
	}
	_ = json.Unmarshal(bodyData, &result)

	info = &FileInfo{
		Filename:      fileName,
		ContentLength: downResp.ContentLength,
		ContentType:   downResp.Header.Get("content-type"),
		Link:          result.Link,
	}

	return
}

func demo1() {
	fileURL := "https://omni-wechat-qa.oss-cn-zhangjiakou.aliyuncs.com/wechat/wxkf/7/wpsduZdQAArDHqPXNhkVZ-03W8uAY87w/wksduZdQAAQeGArm-Q76QSCrHj3orZCw/rpSjWSDDlCWaefl3.mp4"
	info, err := UploadByURL(context.Background(), fileURL, "video", nil)
	if err != nil {
		log.Printf("UploadByURL err: %v\n", err)
		return
	}

	log.Println(info)
}

func Main() {
	demo1()
}

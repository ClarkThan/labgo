package tsoshit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"
)

var (
	// accessKey = os.Getenv("TOS_ACCESS_KEY")
	// secretKey = os.Getenv("TOS_SECRET_KEY")
	accessKey = "AKLTNzU3OGMwNzAwM2ZlNDkyNGIyOWJjYjVlMzliOGVhMjU"
	secretKey = "WlRWaE16RmtPVGcxTjJGa05HRTVNV0UwTmpsaU1HSTFOVE14WXpnMk5ESQ=="
	// Bucket 对于的 Endpoint，以华北2（北京）为例：https://tos-cn-beijing.volces.com
	endpoint = "https://tos-cn-beijing.volces.com"
	region   = "cn-beijing"
	// 填写 BucketName
	bucketName = "photo-meiqia-ai"
	// 填写对象名
	objectKey = "AXYZ/txbCJNHjLw587M9CGmBp.jpg"

	httpClient = &http.Client{}
)

func checkErr(err error) {
	if err != nil {
		if serverErr, ok := err.(*tos.TosServerError); ok {
			fmt.Println("Error:", serverErr.Error())
			fmt.Println("Request ID:", serverErr.RequestID)
			fmt.Println("Response Status Code:", serverErr.StatusCode)
			fmt.Println("Response Header:", serverErr.Header)
			fmt.Println("Response Err Code:", serverErr.Code)
			fmt.Println("Response Err Msg:", serverErr.Message)
		} else if clientErr, ok := err.(*tos.TosClientError); ok {
			fmt.Println("Error:", clientErr.Error())
			fmt.Println("Client Cause Err:", clientErr.Cause.Error())
		} else {
			fmt.Println("Error:", err)
		}
		panic(err)
	}
}

func demo1() {

	// 初始化客户端
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)

	// 上传对象 Body ， 以 string 对象为例
	// body := strings.NewReader("object content")
	fd, err := os.Open("/Users/ranya/Downloads/VywgF3ZIfjSTSUYIi1Hc.jpg")
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()
	// 上传对象
	output, err := client.PutObjectV2(context.Background(), &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: bucketName,
			Key:    objectKey,
		},
		Content: fd,
	})
	checkErr(err)
	fmt.Println("Put Object Request ID: ", output.RequestID)
	fmt.Println("Put Object Response Status Code: ", output.StatusCode)
}

func demo2() {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	// 生成上传对象预签名
	url, err := client.PreSignedURL(&tos.PreSignedURLInput{
		HTTPMethod: enum.HttpMethodPut,
		Bucket:     bucketName,
		Key:        objectKey,
	})
	fmt.Println("url: ", url.SignedUrl)
	for k, v := range url.SignedHeader {
		fmt.Printf(" %s  ->  %s\n", k, v)
	}
	checkErr(err)
	// 上传对象
	body := strings.NewReader("your body reader")
	req, _ := http.NewRequest(http.MethodPut, url.SignedUrl, body)
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("upload object error: %v\n", err)
		return
	}
	io.Copy(os.Stdout, res.Body)
	defer res.Body.Close()
	checkErr(err)
	if res.StatusCode != http.StatusOK {
		panic("unexpected status code")
	}
}

func demo3() {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	// 删除对象
	ret, err := client.DeleteObjectV2(context.Background(), &tos.DeleteObjectV2Input{
		Bucket: bucketName,
		Key:    "Y3QF/txbCJNHjLw587M9CGmBp.jpg", // objectKey,
	})
	if err != nil {
		checkErr(err)
		return
	}

	fmt.Printf("Delete Object ret: %v\n", ret)
}

func Main() {
	demo3()
}

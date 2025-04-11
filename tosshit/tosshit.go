package tsoshit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

// https://www.volcengine.com/docs/6349/132400
func NormalPresigned() {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	// 生成上传对象预签名
	url, err := client.PreSignedURL(&tos.PreSignedURLInput{
		HTTPMethod: enum.HttpMethodPut, // post方法不行啊
		Bucket:     "video-meiqia-ai",
		Key:        "video/20250410/zzzXAtiykBoF/test-hello.mp4",
		Expires:    259200,
	})
	checkErr(err)
	fmt.Printf("headers: %+v\n", url.SignedHeader)
	fmt.Println("url1: ", url.SignedUrl)
	// return
	// 上传对象
	file, err := os.OpenFile("/Users/ranya/Downloads/25_1707037743.mp4", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("open file error: %v\n", err)
	}
	defer file.Close()
	// body := strings.NewReader("your body reader")
	req, err := http.NewRequest(http.MethodPut, url.SignedUrl, file)
	if err != nil {
		log.Fatalf("new request error: %v\n", err)
	}
	res, err := httpClient.Do(req)
	checkErr(err)
	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}
}

// 还是 NormalPresigned 好点儿
func GetSignedPolicyURL() {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	rex, err := client.PreSignedPolicyURL(context.Background(), &tos.PreSingedPolicyURLInput{
		Bucket:              "video-meiqia-ai",
		AlternativeEndpoint: "tos-cn-beijing.volces.com/video/20250410/psZXAtiykBoF/10_25_1707037743.mp4",
		Expires:             259200,
	})
	if err != nil {
		log.Fatalf("get presigned post signature error: %v\n", err)
	}

	fmt.Println("url2: ", rex.GetSignedURLForList(nil))
}

// https://www.volcengine.com/docs/6349/173411
func FormPresigned() {
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	checkErr(err)
	res, err := client.PreSignedPostSignature(context.Background(), &tos.PreSingedPostSignatureInput{
		Bucket:  "video-meiqia-ai",
		Key:     "video/20250410/psZXAtiykBoF/10_25_1707037743.mp4",
		Expires: 259200,
	})
	if err != nil {
		log.Fatalf("get presigned post signature error: %v\n", err)
	}
	fmt.Println("Policy:", res.Policy)
	fmt.Println("OriginPolicy:", res.OriginPolicy)
	fmt.Println("Algorithm:", res.Algorithm)
	fmt.Println("Credential:", res.Credential)
	fmt.Println("Date:", res.Date)
	fmt.Println("Signature:", res.Signature)

	// var x tos.PreSingedPostSignatureOutput
	fmt.Printf("\n\nres: %+v\n", res)
}

func testShit(a, b any) {
	res, err := ScopeEncode(a, b)
	if err != nil {
		fmt.Printf("encode error: %v\n", err)
		return
	}
	fmt.Println(res)
}

func ScopeEncode(privilegeRange, privilegeRangeAgents interface{}) (string, error) {
	switch privilegeRange := privilegeRange.(type) {
	case string:
		if privilegeRange != "self" && privilegeRange != "all" {
			return "", errors.New("invaild privilege range")
		}
		return privilegeRange, nil
	case []any:
		var privilegeGroup []string
		for _, val := range privilegeRange {
			privilegeGroup = append(privilegeGroup, strconv.FormatFloat(val.(float64), 'f', 0, 64))
		}
		mainScope := "groups:" + strings.Join(privilegeGroup, ",")
		subScope, err := subScopeEncode(privilegeRangeAgents)
		if err != nil {
			return "", err
		}
		if mainScope == groupsPrefix {
			return subScope, nil
		}
		return strings.Trim(strings.Join([]string{mainScope, subScope}, "|"), "|"), nil
	case []int:
		var privilegeGroup []string
		for _, val := range privilegeRange {
			privilegeGroup = append(privilegeGroup, fmt.Sprintf("%d", val))
		}
		mainScope := "groups:" + strings.Join(privilegeGroup, ",")
		subScope, err := subScopeEncode(privilegeRangeAgents)
		if err != nil {
			return "", err
		}
		if mainScope == groupsPrefix {
			return subScope, nil
		}
		return strings.Trim(strings.Join([]string{mainScope, subScope}, "|"), "|"), nil
	case nil:
		return subScopeEncode(privilegeRangeAgents)
	}
	return "", errors.New("invaild privilege range")
}

func subScopeEncode(privilegeRangeAgents any) (string, error) {
	switch pVal := privilegeRangeAgents.(type) {
	case nil:
		return "", nil
	case []any:
		if len(pVal) == 0 {
			return "", nil
		}
		var pAgents []string
		for _, val := range pVal {
			pAgents = append(pAgents, strconv.FormatFloat(val.(float64), 'f', 0, 64))
		}
		return "agents:" + strings.Join(pAgents, ","), nil
	case []int:
		if len(pVal) == 0 {
			return "", nil
		}
		var pAgents []string
		for _, val := range pVal {
			pAgents = append(pAgents, fmt.Sprintf("%d", val))
		}
		return "agents:" + strings.Join(pAgents, ","), nil
	}
	return "", errors.New("invaild privilege agents range")
}

const (
	Self         = "self"
	All          = "all"
	groupsPrefix = "groups:"
	agentsPrefix = "agents:"
	separator    = ","
)

type PrivilegeRange struct {
	GroupIDs []int64
	AgentIDs []int64
	Common   string
}

func testDecode(s string) {
	fmt.Printf("got: %+v\n", ScopeDecodeNew(s))
}

// self
// all
// group:1,2,3
// agent:11,22,33
// group:1,2,3|agent:11,22,33
// agent:11,22,33|group:1,2,3
func ScopeDecodeNew(privilegeRange string) PrivilegeRange {
	if privilegeRange == Self || privilegeRange == All {
		return PrivilegeRange{Common: privilegeRange}
	}
	parts := strings.Split(privilegeRange, "|")

	var priv PrivilegeRange
	for _, part := range parts {
		if strings.HasPrefix(part, groupsPrefix) {
			priv.GroupIDs = parseSliceByPrefix(groupsPrefix, part)
		} else if strings.HasPrefix(part, agentsPrefix) {
			priv.AgentIDs = parseSliceByPrefix(agentsPrefix, part)
		}
	}

	if len(priv.AgentIDs) == 0 && len(priv.GroupIDs) == 0 {
		return PrivilegeRange{Common: Self}
	}

	return priv
}

func parseSliceByPrefix(prefix string, str string) []int64 {
	idStr := strings.TrimPrefix(str, prefix)
	if idStr == "" {
		return nil
	}
	idStrs := strings.Split(idStr, separator)
	var ids []int64
	for _, v := range idStrs {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, id)
	}
	return ids
}

func Main() {
	// NormalPresigned()
	// GetSignedPolicyURL()
	// FormPresigned()
	testShit("all", []int{1, 2})
	testShit("self", []int{1, 2})
	testShit("self", 3)
	testShit("self", nil)
	testShit([]int{1, 2, 3}, []int{10, 20, 30})
	testShit(nil, []int{10, 20, 30})
	testShit([]int{1, 2, 3}, nil)
	testShit([]int{0}, []int{})
	testShit([]int{}, []int{111})

	// self
	// all
	// group:1,2,3
	// agent:11,22,33
	// group:1,2,3|agent:11,22,33
	// agent:11,22,33|group:1,2,3
	testDecode("self")
	testDecode("all")
	testDecode("groups:1,2,3")
	testDecode("agents:11,22,33")
	testDecode("groups:1,2,3|agents:11,22,33")
	testDecode("agents:11,22,33|groups:1,2,3")
	testDecode("groups:0|agents:")
}

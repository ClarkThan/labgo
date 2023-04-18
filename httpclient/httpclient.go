package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

var RoundTripper http.RoundTripper

func init() {
	// caCert, err := os.ReadFile("rootCA.crt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	RoundTripper = &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 5,
		DisableCompression:  true,
		IdleConnTimeout:     15 * time.Minute,
		// TLSClientConfig: &tls.Config{
		// 	RootCAs: caCertPool,
		// },
		// #nosec G402
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: RoundTripper,
		// Timeout:   8 * time.Second,
	}
}

type Req struct {
	Question string `json:"question"`
}

type Resp struct {
	Answer string `json:"answer"`
}

type ErrInfo struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Main() {
	client := NewClient()

	q := Req{"分配规则有哪些？"}
	payload, _ := json.Marshal(q)
	url := "http://47.252.6.43:8080/chat"
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()
	// ctx := context.Background()

	var err error

	defer func() {
		if err != nil {
			fmt.Printf("defer err info: %v\n", err)
		}
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	begin := time.Now()
	// resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	resp, err := client.Do(req)
	elapsed := time.Now().Sub(begin)
	// newFunction(ctx)
	fmt.Printf("request elapsed: %s\n", elapsed)

	// data := map[string]any{
	// 	"model":       "text-davinci-003",
	// 	"prompt":      "Say this is a test",
	// 	"max_tokens":  7,
	// 	"temperature": 0,
	// }
	// payload, _ := json.Marshal(data)
	// req, _ := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", bytes.NewReader(payload))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer sk-CjN4fEQzIie5W9gApGljT3BlbkFJstTdZyfBwNxHi41rpAF3PWD")
	// resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		if x, ok := err.(interface{ Timeout() bool }); ok && x.Timeout() {
			fmt.Println("damn it")
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errInfo ErrInfo
		_ = json.NewDecoder(resp.Body).Decode(&errInfo) // nolint: errcheck
		fmt.Printf("err info: %v\n", errInfo)
		return
	}

	var ans Resp
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		fmt.Printf("decode err: %v\n", err)
		return
	}

	fmt.Println("got ans:", ans.Answer)
}

func newFunction(ctx context.Context) {
	fmt.Println(time.Now())
	select {
	case <-ctx.Done():
		fmt.Println("超时:", ctx.Err())
	case <-time.After(3 * time.Second):
		fmt.Println("3s")
	default:
		fmt.Println("没有超时")
	}
	fmt.Println(time.Now())
}

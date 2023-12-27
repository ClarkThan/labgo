package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	RoundTripper http.RoundTripper
	ctx          context.Context = context.Background()
)

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
		// Timeout:   10 * time.Second,
	}
}

type Req struct {
	CorpusID string `json:"corpus_id"`
	Question string `json:"question"`
}

type Answer struct {
	Answer string `json:"answer"`
}

type ErrInfo struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Main() {
	demo4()
}

func demo1() {
	client := NewClient()

	q := Req{CorpusID: "63", Question: "分配规则有哪些？"}
	payload, _ := json.Marshal(q)
	url := "http://47.252.6.43:8080/chat"
	// ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	// defer cancelFn()
	ctx := context.Background()

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
		var ei interface{ Timeout() bool }
		if errors.As(err, &ei) && ei.Timeout() {
			fmt.Println("damn it, fuck you")
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

	var ans Answer
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		fmt.Printf("decode err: %v\n", err)
		return
	}

	fmt.Println("got ans:", ans.Answer)
}

func demo2() {
	client := NewClient()

	q := Req{CorpusID: "63", Question: "分配规则有哪些？"}
	payload, _ := json.Marshal(q)
	url := "http://47.252.6.43:8080/chat"
	start := time.Now()
	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	elapsed := time.Since(start).Milliseconds()
	if err != nil {
		var e interface{ Timeout() bool }
		if errors.As(err, &e) && e.Timeout() {
			fmt.Printf("超时拉: %dms", elapsed)
		}
		return
	}

	defer resp.Body.Close()
	var ans Answer
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		fmt.Printf("decode gptbot resp failed;  %v", err)
		return
	}

	fmt.Printf("(%dms)答案: %s", elapsed, ans.Answer)
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

type Demo3Req struct {
	Question string `json:"question"`
}

type CodeEnum string

const (
	Yes CodeEnum = "Y"
	No  CodeEnum = "N"
)

type AnswerTypeEnum uint8

const (
	Ansr AnswerTypeEnum = iota + 1
	List
	Empt
)

type AnswerItem struct {
	AnswerType       AnswerTypeEnum `json:"answer_type"`
	AnswerText       string         `json:"answer"`
	StandardQuestion string         `json:"standard_question"`
	// RecommendAnswerList []string       `json:"recommend_answer_list"`

}

type Demo3Resp struct {
	Code CodeEnum     `json:"code"`
	Body []AnswerItem `json:"body"`
	// Body json.RawMessage `json:"body"`
}

type MessageResp struct {
	Understood bool      `json:"understood"`
	Contents   []Content `json:"contents"`
	HitType    string    `json:"hit_type"`
}

type Content struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"` // 输入 string ｜ 输出 SearchKnowResponse
}

func demo3() {
	client := NewClient()
	// q := Demo3Req{Question: "分单和并单的作用是什么"}
	// q := Demo3Req{Question: "我是谁"}
	// q := Demo3Req{Question: ""}
	q := Demo3Req{Question: "如何查询集装箱超期使用费？"}
	payload, _ := json.Marshal(q)
	// url := "https://gw-api-hk-di1.sit.cmft.com/ics-hk/chat/workbench?tenant_code=TEST01"
	url := "https://ics-hk.cm-worklink.com/ics-hk/chat/workbench?tenant_code=JYCC"
	start := time.Now()
	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Panicf("请求失败: %v\n", err)
	}
	elapsed := time.Since(start).Milliseconds()

	defer resp.Body.Close()
	var ans Demo3Resp
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		fmt.Printf("decode resp failed;  %v", err)
		return
	}

	fmt.Printf("耗时 (%dms)", elapsed)
	// fmt.Printf("耗时 (%dms) 答案: %+v\n", elapsed, ans)
	// expected := Demo3Resp{
	// 	Code: "Y",
	// 	Body: []AnswerItem{
	// 		{
	// 			AnswerText:       "分单和并单的作用是对订单的提单号进行管理和调整，方便客户进行货物的跟踪和管理，同时也可以对拼单进行管理和调整，实现更加灵活的订舱操作。",
	// 			AnswerType:       Ansr,
	// 			StandardQuestion: "分单和并单的作用是什么？",
	// 		},
	// 	},
	// }
	// expected := Demo3Resp{
	// 	Code: "Y",
	// 	Body: []AnswerItem{
	// 		{
	// 			AnswerText:       "",
	// 			AnswerType:       List,
	// 			StandardQuestion: []any{},
	// 		},
	// 	},
	// }

	// if !cmp.Equal(ans, expected) {
	// 	log.Printf("damn: %s\n", cmp.Diff(ans, expected))
	// } else {
	// 	log.Println("no problem")
	// }

	understood := false
	content := []Content{}
	if ans.Code == Yes && len(ans.Body) > 0 {
		for _, a := range ans.Body {
			if a.AnswerType == Ansr {
				understood = true
				content = append(content, Content{
					Type: "rich_text",
					Body: map[string]any{"answer": a.AnswerText, "answer_rich": fmt.Sprintf("<p>%s</p>", a.AnswerText), "question": a.StandardQuestion},
				})
			}
		}
	}

	// if ans.Code == Yes {
	// 	var ansItems []AnswerItem
	// 	json.Unmarshal(ans.Body, &ansItems)
	// 	for _, a := range ansItems {
	// 		if a.AnswerType == Ansr {
	// 			understood = true
	// 			content = append(content, Content{
	// 				Type: "rich_text",
	// 				Body: map[string]any{"answer": a.AnswerText, "answer_rich": fmt.Sprintf("<p>%s</p>", a.AnswerText), "question": a.StandardQuestion},
	// 			})
	// 		}
	// 	}
	// }

	data := &MessageResp{
		Understood: understood,
		HitType:    "zhaoshang",
		Contents:   content,
	}

	dat, _ := json.Marshal(data)
	log.Println("----> ", string(dat))
}

// --------------------- demo4 -----------------------

type Session struct {
	HTTPClient *http.Client
}

func NewSession() *Session {
	return &Session{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Session) MakeRequest(ctx context.Context, endpoint string, input, output any) error {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "hikari-gptbot")

	resp, err := s.HTTPClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("error making wenxin yiyan request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 || bytes.Contains(respBody, []byte(`"error_code":`)) {
		return &APIError{
			// StatusCode: resp.StatusCode,
			Payload: respBody,
		}
	}

	return json.Unmarshal(respBody, output)
}

// APIError is returned from API requests if the API responds with an error.
type APIError struct {
	// StatusCode int
	Payload []byte
}

func (e *APIError) Error() string {
	// return fmt.Sprintf("status_code=%d, payload=%s", e.StatusCode, e.Payload)
	return string(e.Payload)
}

type Question struct {
	CorpusID string `json:"corpus_id"`
	Question string `json:"question"`
}

func demo4() {
	question := Question{Question: "讲一个笑话", CorpusID: "23"}
	payload, _ := json.Marshal(question) // nolint: errcheck
	client := NewClient()
	resp, err := client.Post("http://127.0.0.1:8080/chat", "application/json", bytes.NewReader(payload))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("damn has error")
		var e interface{ Timeout() bool }
		if errors.As(err, &e) && e.Timeout() {
			log.Println("api timeout")
		} else {
			log.Printf("got error: %v\n", err)
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("not 200", resp.StatusCode)
		// return
	}

	var ans Answer
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		return
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read resp body failed: %v\n", err)
	}
	log.Println("resp body: ", string(dat))

	log.Println("got answer: ", ans.Answer)
}

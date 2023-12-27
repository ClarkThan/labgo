package styleshit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type YMCA interface {
	Shout() string
}

func UseYMCA(y YMCA) {
	fmt.Println("i got ", y.Shout())
}

type Impl struct {
	Poetry string
}

func (i *Impl) Shout() string {
	return i.Poetry
}

// func RetYMCA(s string) *Impl {
func RetYMCA(s string) YMCA {
	return &Impl{Poetry: s}
}

func demo1() {
	y := RetYMCA("god damn it!")
	UseYMCA(y)
}

type CommonResp struct {
	Name  string
	Value any
}

// type Anser interface {
// 	GetAns() string
// }

type Resp1 struct {
	Answer string
}

// func (r *Resp1) GetAns() string {
// 	return r.Answer
// }

type Resp2 struct {
	Question string
	Answer   string
}

// func (r *Resp2) GetAns() string {
// 	return r.Answer
// }

func ret_resp1() CommonResp {
	return CommonResp{Name: "resp1", Value: Resp1{Answer: "resp1_answer"}}
	// return CommonResp{Name: "resp1", Value: map[string]string{"answer": "resp1_answer"}}
}

func ret_resp2() CommonResp {
	return CommonResp{Name: "resp1", Value: Resp2{Question: "resp2_question", Answer: "resp2_answer"}}
}

type InterSet struct {
	Answer string `json:"answer"`
}

func demo2() {
	c1 := ret_resp1()
	c2 := ret_resp2()
	// c1..Answer + c2..Answer
	ans1 := reflect.ValueOf(c1.Value).FieldByName("Answer")
	ans2 := reflect.ValueOf(c2.Value).FieldByName("Answer")
	log.Println(ans1, ans2)
	c1Dat, _ := json.Marshal(c1.Value)
	var c1tmp InterSet
	_ = json.Unmarshal(c1Dat, &c1tmp)
	c2Dat, _ := json.Marshal(c2.Value)
	var c2tmp InterSet
	_ = json.Unmarshal(c2Dat, &c2tmp)

	log.Println(c1tmp.Answer, c2tmp.Answer)

	ret := InterSet{Answer: c1tmp.Answer + c2tmp.Answer}
	retDat, _ := json.Marshal(ret)
	log.Println(string(retDat))
}

func Main() {
	demo2()
}

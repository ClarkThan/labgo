package reflectshit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var (
	cols = []string{"var1", "var2", "var3", "var4", "var5"}
)

type ManyCol struct {
	Var1  string `json:"var1,omitempty" db:"var1"`
	Var2  string `json:"var2,omitempty" db:"var2"`
	Var3  string `json:"var3,omitempty" db:"var3"`
	Var4  string `json:"var4,omitempty" db:"var4"`
	Var5  string `json:"var5,omitempty" db:"var5"`
	Text1 string `json:"text1,omitempty" db:"text1"`
	Text2 string `json:"text2,omitempty" db:"text2"`
	Text3 string `json:"text3,omitempty" db:"text3"`
	Text4 string `json:"text4,omitempty" db:"text4"`
	Text5 string `json:"text5,omitempty" db:"text5"`
}

func (m *ManyCol) Normal() {
	if m.Var1 == "" {
		m.Text1 = ""
	} else if m.Text1 == "" {
		m.Text1 = fmt.Sprintf("请拨打 [%s]？", m.Var1)
	}
}

func Main() {
	demo2()
}

func demo1() {
	m := ManyCol{
		Var1: "tel",
	}
	m.Normal()

	field := "var1"
	textField := strings.Replace(field, "var", "Text", 1)
	text1 := reflect.ValueOf(m).FieldByName(textField)
	fmt.Println(text1)
}

type AgentSt struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Age    int8   `json:"age"`
}

type Event struct {
	Agent   AgentSt `json:"agent"`
	TrackID string  `json:"track_id"`
}

func helper(evt any) {
	// field := reflect.ValueOf(evt).FieldByName("Agent")
	// field, ok := reflect.TypeOf(evt).Elem().FieldByName("Agent")
	// if ok {
	// 	fmt.Println(field.Name)
	// }
	x, ok := reflect.TypeOf(evt).FieldByName("Agent")
	if !ok {
		log.Fatalf("damn")
	}
	// fmt.Println("over 1")
	x.Tag = reflect.StructTag(`json:"-"`)
	// fmt.Println("over 2")
	dat, _ := json.Marshal(evt)
	fmt.Println("changed:", string(dat))
}

func demo2() {
	evt := Event{
		TrackID: "fuckyou",
		Agent: AgentSt{
			Name:   "air",
			Age:    23,
			Avatar: "http://meiqia.com",
		},
	}
	dat, _ := json.Marshal(evt)
	fmt.Println("origin:", string(dat))
	helper(evt)
}

package reflectshit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	// "github.com/fatih/structs"
	"github.com/RussellLuo/structs"
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

func demo3() {
	// var n int8 = 23
	// var m M
	var num int16 = 123
	m := M{Name: "Air", Age: 23, ok: inner{Foo: "hello", Bar: []*int16{&num}}}
	dict := struct2Map(m)

	fmt.Println(reflect.TypeOf(dict["Ok"]).Kind())

	d1, ok1 := dict["Ok"].(inner)
	if !ok1 {
		fmt.Println("damn 1")
		return
	}

	fmt.Println(d1.Body)

	// fmt.Println("dict", dict)
	dat, err := json.Marshal(dict)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("dict", string(dat))

	sts := structs.New(m)
	sts.TagName = "json"
	st := sts.Map()
	dat, err = json.Marshal(st)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("stucts", string(dat))
}

type inner struct {
	Foo  string             `json:"foo"`
	Bar  []*int16           `json:"bar"`
	Body map[string]*string `json:"body"`
}

type M struct {
	Name string   `json:"name"`
	Age  uint8    `json:"age"`
	Fuck string   `json:"fuck,omitempty"`
	Addr []string `json:"addr"`
	Node inner    `json:"node"`
	Ptr1 *inner   `json:"ptr1,omitempty,omitnested"`
	Ptr2 ***inner
	ok   inner `json:"ok"`
}

func demo4() {
	val := "bar"
	var num int16 = 123
	in1 := inner{Foo: "damn shit"}
	in2 := &in1
	in3 := &in2
	m := M{Name: "Air", Age: 23, ok: inner{Foo: "hello", Bar: []*int16{&num}}, Node: inner{Foo: "node", Bar: []*int16{&num}, Body: map[string]*string{"foo": &val}}, Ptr2: &in3}
	dict := struct2Map(&m)
	fmt.Printf("dict: \n%#v\n\n", dict)

	fmt.Println("dict avg alloc:", testing.AllocsPerRun(50, func() {
		_ = struct2Map(&m)
	}))
	// fmt.Println(dict["ok"].(map[string]any)["bar"].([]any))
	// dat, err := json.Marshal(dict)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("json:", string(dat))

	st := structs.New(m)
	st.TagName = "json"
	mmap := st.Map()
	fmt.Printf("mmap: \n%#v\n", mmap)
	fmt.Println("mmap avg alloc:", testing.AllocsPerRun(50, func() {
		_ = structs.New(m).Map()
	}))
	// fmt.Println(mmap["ok"].(map[string]any)["bar"].([]int16))

	// s := []int8{11, 22, 33}
	// rv := reflect.ValueOf(s)
	// fmt.Println(rv.Kind(), rv.Len())
	// fmt.Println(rv.Index(0))

	// m := map[string]int{"foo": 1, "bar": 2}
	// rv := reflect.ValueOf(m)
	// fmt.Println(rv.Kind(), rv.Len())
	// iter := rv.MapRange()
	// for iter.Next() {
	// 	fmt.Println(iter.Key(), "->", iter.Value())
	// }

	// var n int16 = 23
	// var ptr *int16
	// rv := reflect.ValueOf(ptr)
	// fmt.Println(rv.Kind() == reflect.Pointer, rv.IsZero(), rv.IsNil())
	// fmt.Println(rv.Elem())
	// fmt.Println(rv.Kind(), rv.Elem().Interface())
}

type Demo5 struct {
	Name   string         `json:"name"`
	Age    int8           `json:"age"`
	Ok     bool           `json:"ok"`
	Score  float64        `json:"score"`
	Meters []float32      `json:"meters,omitnested"`
	Info   map[string]any `json:"info,omitnested"`
	Body   InnerBody      `json:"body,omitnested"`
}

type InnerBody struct {
	Nickname string `json:"nickname"`
	Height   int    `json:"heigh"`
}

func demo5() {
	m := Demo5{
		Name:   "jordan",
		Age:    23,
		Ok:     true,
		Score:  123.45,
		Meters: []float32{12.1, 34.1, 56.1},
		Info:   map[string]any{"addr": "Sichuan", "no": 45},
		Body:   InnerBody{Nickname: "air", Height: 100},
	}

	fmt.Println("struct avg alloc:", testing.AllocsPerRun(50, func() {
		st := structs.New(m)
		st.TagName = "json"
		_ = st.Map()
	}))

	fmt.Println("myself avg alloc:", testing.AllocsPerRun(50, func() {
		_ = struct2Map(&m)
	}))

	st := structs.New(m)
	st.TagName = "json"
	s1 := st.Map()
	fmt.Printf("%#v\n\n", s1)

	s2 := struct2Map(m)
	fmt.Printf("%#v\n", s2)
}

func demo6() {
	type B struct {
		Time []time.Time `structs:"time"`
	}
	type A struct {
		Time time.Time `structs:"time"`
		B    B         `structs:"b"`
	}

	f := func(value reflect.Value) (interface{}, error) {
		switch v := value.Interface().(type) {
		case time.Time:
			return v.Format(time.RFC3339), nil
		}
		return value.Interface(), nil
	}

	date := time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC)
	s := structs.New(A{Time: date, B: B{Time: []time.Time{date}}})
	s.EncodeHook = f

	got := s.Map()
	// want := map[string]interface{}{
	// 	"time": "2021-08-15T00:00:00Z",
	// 	"b": map[string]interface{}{
	// 		"time": "2021-08-15T00:00:00Z",
	// 	},
	// }

	fmt.Printf("%#v\n", got)
	// fmt.Printf("%#v\n", want)
}

type EffectiveSetting struct {
	Condition string `json:"condition"`
	Rule      string `json:"rule"`
	Count     int64  `json:"count"`
}

type RuleSettings struct {
	EffectiveSetting EffectiveSetting `json:"effective_setting"`
	EffectiveTime    string           `json:"effective_time"`
	OK               bool             `json:"ok"`
}

func IsZero(v any) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func demo7() {
	var r RuleSettings
	r.OK = true
	// rv := reflect.ValueOf(r)
	fmt.Println(IsZero(r))
	var x map[string]struct{}
	fmt.Println(IsZero(x))
	y := make(map[string]struct{})
	fmt.Println(IsZero(y))
	z := []string{}
	fmt.Println(IsZero(z))
}

func Main() {
	_, _ = demo8(0)
	_, _ = demo8(9)
	_, _ = demo8(3)
}

func demo8(n int) (ret int, err error) {
	defer func() {
		fmt.Printf("ret: %d, err: %v\n", ret, err)
	}()
	if n == 0 {
		return 0, fmt.Errorf("zero")
	}

	if n > 3 {
		ret = 33
		return
	}

	return n, nil
}

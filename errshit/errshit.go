package errshit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type X struct {
	Name string
}

func (x X) Error() string {
	return fmt.Sprintf("X Err: %s", x.Name)
}

func (x X) Shit() string {
	return "no " + x.Name
}

func demo1() {
	var e1 error = X{"foo"}
	e2 := fmt.Errorf("Wrapper2  %w", e1)
	e3 := fmt.Errorf("Wrapper3  %w", e2)
	fmt.Println(e3)

	var eo error = e3
	for {
		tmp := errors.Unwrap(eo)
		if tmp == nil {
			break
		}
		eo = tmp
	}

	fmt.Println(eo)
	// if ee3, ok := eo.(interface{ Shit() string }); ok {
	// 	fmt.Println(ee3.Shit())
	// }
	var perr interface{ Shit() string }
	if errors.As(eo, &perr) {
		fmt.Println("-----", perr.Shit())
	}
}

func copyMap(m map[string]any) map[string]any {
	ret := make(map[string]any)
	for k, v := range m {
		ret[k] = v
	}

	return ret
}

func demo2() {
	m1 := map[string]any{
		"foo": "bar",
		"info": map[string]any{
			"name":    "mj",
			"numbers": []int64{9, 12, 23},
			"age":     60,
		},
		"baz": 100,
	}

	m1Bs, _ := json.Marshal(m1)
	log.Println(string(m1Bs))

	m2 := copyMap(m1)
	m2Info, _ := m2["info"].(map[string]any)
	m2InfoNumbers, _ := m2Info["numbers"].([]int64)
	m2Info["numbers"] = append(m2InfoNumbers, 45)
	m2["info"] = m2Info

	m1Bss, _ := json.Marshal(m1)
	log.Println(string(m1Bss))
}

func Main() {
	demo2()
}

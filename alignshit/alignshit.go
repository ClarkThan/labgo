package alignshit

import (
	"fmt"
	"unsafe"
)

type TimeSpec struct {
	Type  string `json:"type"`
	Begin string `json:"begin"`
	End   string `json:"end"`
	Days  []int  `json:"days"`
}

func Main() {
	var t TimeSpec
	fmt.Println(unsafe.Sizeof(t))
}

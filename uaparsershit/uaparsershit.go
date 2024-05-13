package uaparsershit

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/ua-parser/uap-go/uaparser"
)

var (
	parser = uaparser.NewFromSaved()

	huanweiRegex = regexp2.MustCompile(`\((?:(\w+;\s?)?)(?<os>OpenHarmony|HarmonyOS)\s+(?<os_ver>\d+(?:\.\d+)+)\)`, 0)
)

func demo1() {
	s := "Mozilla/5.0 (Phone; OpenHarmony 4.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 ArkWeb/4.1.6.1 Mobile HuaweiBrowser/1.2.7.343"
	// s = "Mozilla/5.0 (Linux; Android 10; HarmonyOS; ANA-AN00; HMSCore 6.11.0.302) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 HuaweiBrowser/13.0.6.302 Mobile Safari/537.36"
	// s = "Mozilla/5.0 (Linux; Android 12; ANG-AN00 Build/HUAWEIANG-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/97.0.4692.98 Mobile Safari/537.36 T7/13.34 BDOS/1.0 (HarmonyOS 3.0.0) SP-engine/2.72.0 baiduboxapp/13.34.5.10 (Baidu; P1 12) NABar/1.0"
	ua := parser.Parse(s)
	fmt.Printf("%#v\n", ua.Device)
	fmt.Printf("%#v\n", ua.Os)
	fmt.Printf("%#v\n", ua.UserAgent)
	// fmt.Printf("%#v\n", ua.Os)
	// fmt.Printf("%#v\n", ua.Os)

	os, osVer := checkForHuawei(ua, s)
	fmt.Println(os, osVer)
}

func checkForHuawei(uc *uaparser.Client, userAgent string) (os, osVer string) {
	if uc.Device.Brand == "Huawei" && uc.Os.Family != "Android" {
		fmt.Println("enter huawei")
		matched, err := huanweiRegex.FindStringMatch(userAgent)
		if err != nil || matched == nil {
			return
		}
		os := matched.GroupByName("os").Capture.String()
		osVer := matched.GroupByName("os_ver").Capture.String()
		return os, osVer
	}

	return
}

func Main() {
	demo1()
}

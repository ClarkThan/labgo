package regexpshit

import (
	"fmt"
	"log"
	"regexp"
	"regexp/syntax"
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	UnicodeRegex       = regexp.MustCompile(`\\u[0-9a-f]{4}`)
	MetaRegexExtractor = regexp2.MustCompile(`(\/)?(?<pattern>.*)(?(1)\/)(?<flags>[gismuy]*)`, 0)
	SerialNORegex      = regexp.MustCompile(`\d+`)
	// PhoneRegex         = regexp2.MustCompile(`(?<!\d)(?<phone>1[3|4|5|6|7|8|9]\d{9})(?!\d)`, 0)
	PhoneRegex    = regexp2.MustCompile(`(?<!\d)(?<phone>1[3|4|5|6|7|8|9]\d{2}\s?\d{4}\s?\d{3})(?!\d)|(?<!\d)(?<phone>1[3|4|5|6|7|8|9]\d{1}\s?\d{4}\s?\d{4})(?!\d)`, 0)
	mobileRegexp  = regexp2.MustCompile(`(?<!\d)(1[3|4|5|6|7|8|9]\d{2}\s?\d{4}\s?\d{3})(?!\d)|(?<!\d)(1[3|4|5|6|7|8|9]\d{1}\s?\d{4}\s?\d{4})(?!\d)`, 0)
	mobileRegexp1 = regexp2.MustCompile(`(?<!\d)(1[3|4|5|6|7|8|9]\d{2}(?<sep>[-\s/]?)\d{4}\<sep>\d{3})(?!\d)|(?<!\d)(1[3|4|5|6|7|8|9]\d{1}(?<sep>[-\s/]?)\d{4}\<sep>\d{4})(?!\d)`, 0)
	mobileRegexp2 = regexp2.MustCompile(`(?<!\d)(?<phone>1[3|4|5|6|7|8|9]\d{2}(?<sep>[ -/]?)?\d{4}\<sep>\d{3})(?!\d)|(?<!\d)(?<phone>1[3|4|5|6|7|8|9]\d{1}(?<sep>[ -/]?)\d{4}\<sep>\d{4})(?!\d)`, 0)
	reTel         = regexp2.MustCompile(`(?<!\d)(1[3|4|5|6|7|8|9]\d{2}\s?\d{4}\s?\d{3}|1[3|4|5|6|7|8|9]\d{1}\s?\d{4}\s?\d{4}|0(?:\d{2}-\d{8}|\d{3}-\d{7})|\d{3}-\d{3}-\d{4})(?!\d)`, 0)
	reBadTel      = regexp2.MustCompile(`(?<!\d)(1[3|4|5|6|7|8|9]\d{2}\s?\d{4}\s?\d{2,4}|1[3|4|5|6|7|8|9]\d{1}\s?\d{4}\s?\d{3,5})(?!\d)`, 0)
	reSep         = regexp.MustCompile(`[-\s/]`)
)

func demo1() {
	raw := `[\u4e00-\u9fa5]{2,8}`
	fmt.Printf("raw = %q\n", raw)
	_, e := regexp.Compile(raw)
	if e != nil && strings.Contains(e.Error(), string(syntax.ErrInvalidEscape)) {
		// [\u4e00-\u9fa5]{2,8} -> [\x{4e00}-\x{9fa5}]{2,8}
		raw = UnicodeRegex.ReplaceAllStringFunc(raw, func(rs string) string {
			return fmt.Sprintf("\\x{%s}", strings.TrimLeft(rs, `\u`))
		})
		// fmt.Println(raw)
		fmt.Printf("raw = %q\n", raw)
		matched := regexp.MustCompile(raw).FindString("我是中国人 Bruce Lee.")
		fmt.Println(matched)
	}
}

func demo2() {
	MetaRegexExtractor.FindStringMatch("/^(?!SAT|TOEFL|IELTS).*/")
	r := regexp2.MustCompile(`^(?!SAT|TOEFL|IELTS).*`, 0)
	fmt.Println(r)
	m, err := r.FindStringMatch("sAT123")
	if err != nil || m == nil {
		log.Fatal("shit match")
	}

	// for _, g := range m.Groups() {
	// 	log.Println(g.String())
	// }

	log.Println(m.Group.String())
}

func demo3() {
	_, target := MatchByCustomRegexp(`\b(\d{4} \d{3} \d{4})|(\d{3} \d{4} \d{4})\b`, "号码是 1314 553 4313")
	log.Println(target == "1314 553 4313")
	log.Println("regexp matched:", target)
}

func MatchByCustomRegexp(validationStr, content string) (bool, string) {
	rawRegex := strings.Trim(validationStr, " ")
	if rawRegex == "" {
		return false, ""
	}

	// 举个栗子: "/[0-9a-z]{3,5}/is"  ->  pattern: [0-9a-z]{3,5}    flags: is
	matched, err := MetaRegexExtractor.FindStringMatch(rawRegex)
	if err != nil || matched == nil {
		return false, ""
	}

	// gps := matched.Groups()
	pat := matched.GroupByName("pattern").Capture.String()
	log.Println("----------", pat)
	flags := matched.GroupByName("flags").Capture.String()
	var fixedFlags, pattern string
	// JavaScript 支持 gismuy, GoLang 支持 ismU, 这里过滤一下
	for _, flag := range []string{"i", "s", "m"} {
		if strings.Contains(flags, flag) {
			fixedFlags += flag
		}
	}
	if fixedFlags == "" {
		pattern = pat
	} else {
		pattern = fmt.Sprintf("(?%s)%s", fixedFlags, pat)
	}

	// 将用户自定义中可能包含的 \uxxxx 转化为 \x{xxxx}, 比如：\u4e00 -> \x{4e00}
	pattern = UnicodeRegex.ReplaceAllStringFunc(pattern, func(us string) string {
		return fmt.Sprintf(`\x{%s}`, strings.TrimLeft(us, `\u`))
	})

	// target := regexp.MustCompile(pattern).FindString(content)
	// return target != "", target

	log.Println("---------->>>>>", pattern)
	r2e, err := regexp2.Compile(pattern, 0)
	if err != nil || r2e == nil {
		return false, ""
	}

	r2m, err := r2e.FindStringMatch(content)
	if err != nil || r2m == nil {
		return false, ""
	}

	target := r2m.String()
	return target != "", target
}

func demo4() {
	log.Printf("nos = %v\n", SerialNORegex.FindAllStringSubmatch("机器人客服-好+23-45", -1))
}

func MatchPhone(content string) (bool, string) {
	matched, err := PhoneRegex.FindStringMatch(content)
	if err != nil || matched == nil {
		return false, ""
	}

	return true, matched.GroupByName("phone").Capture.String()
}

func demo5() {
	ok, val := MatchPhone("19182255030")
	ok, val = MatchPhone("191 8225 5030")
	// ok, val = MatchPhone("1918 2255 030")
	if ok {
		fmt.Println(val)
	}
}

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func extractPhoneByGroup(s string) string {
	matched, err := mobileRegexp2.FindStringMatch(s)
	if err != nil || matched == nil {
		return ""
	}

	phone := matched.GroupByName("phone").Capture.String()
	return reSep.ReplaceAllString(phone, "")
}

func demo6() {
	vals := regexp2FindAllString(mobileRegexp, "191 8225 5030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "191-8225-5030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "1918-2255-030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "19182255030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "191 8225 5030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "1918 2255 030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "191-8225 5030")
	fmt.Println(vals)
	fmt.Println("-----------------------------")
	vals = regexp2FindAllString(mobileRegexp1, "191/8225/5030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "1918/2255/030")
	fmt.Println(vals)
	vals = regexp2FindAllString(mobileRegexp1, "191/8225 5030")
	fmt.Println(vals)
}

func demo6_1() {
	vals := extractPhoneByGroup("这个 191-8225-5030，我的号码")
	fmt.Println(vals)
	vals = extractPhoneByGroup("这个1918-2255-030,我的号码")
	fmt.Println(vals)
	vals = extractPhoneByGroup("这个19182255030 我的号码")
	fmt.Println(vals)
	vals = extractPhoneByGroup("这个 191 8225 5030")
	fmt.Println(vals)
	vals = extractPhoneByGroup("这个  1918 2255 030我的号码")
	fmt.Println(vals)
	fmt.Println("---------------")
	vals = extractPhoneByGroup("我的号码191/8225/5030 就这个")
	fmt.Println(vals)
	vals = extractPhoneByGroup("我的号码191/8225/5030 就这个")
	fmt.Println(vals)
	vals = extractPhoneByGroup("我的号码191/8225/5030 就这个")
	fmt.Println(vals)
	fmt.Println(reSep.ReplaceAllString("hello 	world", ""))
	fmt.Println(reSep.ReplaceAllString("hello world", ""))
	fmt.Println(reSep.ReplaceAllString("hello	world", ""))
}

func demo7() {
	vals := regexp2FindAllString(reTel, "1918-2255-0311")
	fmt.Println(vals)
	vals = regexp2FindAllString(reBadTel, "1918 2255 0311")
	// vals := regexp2FindAllString(reTel, "2601692192")
	// vals := regexp2FindAllString(reTel, "330360719@qq.com")
	fmt.Println(vals)
}

func Main() {
	demo6_1()
}

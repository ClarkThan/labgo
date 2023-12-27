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
)

func Main() {
	demo4()
}

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

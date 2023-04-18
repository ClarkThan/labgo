package regexpshit

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
)

var (
	UnicodeRegex = regexp.MustCompile(`\\u[0-9a-f]{4}`)
)

func Main() {
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

package templateshit

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/valyala/fasttemplate"
)

type Obj struct {
	BucketName string
	ObjectName string
}

func demo1() {
	templateStr := "https://{{.BucketName}}.meiqiausercontent.com/{{.ObjectName}}"
	tpl, _ := template.New("media").Parse(templateStr)
	// obj := Obj{BucketName: "pics.meiqia.com"} //, ObjectName: "xeqwrefqdqcaf34"}
	obj := map[string]any{"BucketName": "pics.meiqia.com"}

	// var b bytes.Buffer
	var b strings.Builder
	// w := bufio.NewWriter(&b)
	err := tpl.Execute(&b, obj)
	// err := tpl.Execute(&b, map[string]string{"BucketName": "fuckyou"})
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Println(b.String())

	str := "hello world"

	// Convert string to MD5 hash
	hash := md5.Sum([]byte(str))
	m := md5.New()
	m.Write([]byte(str))
	hash1 := m.Sum(nil)
	// Encode MD5 hash to Base16 (hexadecimal)
	hexStr := hex.EncodeToString(hash[:])

	fmt.Println("Original String:", str)
	fmt.Println("MD5 Hash :", hash)
	fmt.Println("MD5 Hash1:", hash1)
	fmt.Println("Base16 Encoded:", hexStr)
}

func demo2() {
	templateStr := "https://{{.BucketName}}.meiqiausercontent.com/{{.ObjectName}}\n"
	t := fasttemplate.New(templateStr, "{{", "}}")
	t.Execute(os.Stdout, map[string]interface{}{
		"BucketName": "laigu",
		"age":        "18",
	})

	fmt.Println()

	t.ExecuteFunc(os.Stdout, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "name":
			return w.Write([]byte("hjw"))
		case "age":
			return w.Write([]byte("20"))
		}

		return 0, nil
	})
}

func Main() {
	demo1()
}

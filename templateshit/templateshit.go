package templateshit

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"text/template"
)

type Obj struct {
	BucketName string
	ObjectName string
}

func Main() {
	templateStr := "https://{{.BucketName}}.meiqiausercontent.com/{{.ObjectName}}"
	tpl, _ := template.New("media").Parse(templateStr)
	obj := Obj{BucketName: "pics.meiqia.com", ObjectName: "xeqwrefqdqcaf34"}

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
	// Encode MD5 hash to Base16 (hexadecimal)
	hexStr := hex.EncodeToString(hash[:])

	fmt.Println("Original String:", str)
	fmt.Println("MD5 Hash:", hash)
	fmt.Println("Base16 Encoded:", hexStr)
}

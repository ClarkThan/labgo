package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/minio/sha256-simd"
)

func BytesToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

func String2Bytes(s string) []byte {
	ss := (*[2]uintptr)(unsafe.Pointer(&s))
	bs := [3]uintptr{ss[0], ss[1], ss[1]}
	return *(*[]byte)(unsafe.Pointer(&bs))
}

func StringToBytesV0(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		int
	}{s, len(s)}))
}

func StringToBytesV1(s string) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}

func BytesToStringV1(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	runtime.KeepAlive(&bytes) // this line is essential.
	return s
}

func BytesToStringV2(bytes []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}))
}

func s2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

func Curl(uri string) {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	host := u.Hostname()
	port := u.Port()
	path := u.Path

	if port == "" {
		port = "80"
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Fprintf(conn, "GET %s HTTP/1.0\r\nHost: %s\r\n\r\n", path, host)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf[:n]))
}

func foo(n int) (int, error) {
	if n > 100 {
		return 0, errors.New("too big")
	}

	return n, nil
}

func bar(n int) {
	fmt.Println(n)
	_, _ = foo(n)
}

// 循环右移
func rightRotate(num, k, bits int) string {
	k = k % bits // 确保 k 在 [0, bits) 范围内
	val := (num >> k) | (num << (bits - k))
	return fmt.Sprintf("0x%x", val)
}

// rightRotate(0x23, 8, 32)  ==  0x23000000
// rightRotate(0xf2, 24, 32)  ==  0xf200

func Verify(passwd, dbPwd, dbSalt string) error {
	inputPwd := sha256.Sum256([]byte(passwd + dbSalt))
	if hex.EncodeToString(inputPwd[:]) != dbPwd {
		return errors.New("bad pwd")
	}
	return nil
	//return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
}

// Verify("12345678A", "9cbd66f282f4b7347bc065d26cb3ac6ba7756d14bd570db9a2dbc14db92d2e06", "d74f1e50fb8c63fbc67cd1a47cdfd38c")

func GoID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	return strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
}

func TraceMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// fmt.Printf("stats: %#v\n", m)
	fmt.Printf("Alloc = %v HeapAlloc = %v TotalAlloc = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024/1024, m.HeapAlloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
}

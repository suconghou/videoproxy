package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// Log print to stdout
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	t   = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
)

// JSONPut resp json,如果v是byte类型,我们应直接使用,byte类型再json.Marshal就是base64字符串了,string类型经json.Marshal后转为对于的byte
func JSONPut(w http.ResponseWriter, v interface{}, status int, age int) (int, error) {
	var (
		bs  []byte
		err error
	)
	if bb, ok := v.([]byte); !ok {
		if bs, err = json.Marshal(v); err != nil {
			return 0, err
		}
	} else {
		bs = bb
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Max-Age", "864000")
	h.Set("Cache-Control", fmt.Sprintf("public,max-age=%d", age))
	w.WriteHeader(status)
	return w.Write(bs)
}

// GzipEncode gzip data
func GzipEncode(data []byte) ([]byte, error) {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

// DecodeVid video id
func DecodeVid(str string, r1 int, r2 int) (string, error) {
	if len(str) < 1 {
		return "", fmt.Errorf("密文不合规")
	}
	var l = len(t)
	var base = []byte(t)
	var bytestr = []byte(str)
	var n = 0
	for _, char := range bytestr[0 : len(bytestr)-1] {
		n += int(char)
	}
	if base[(n+r2)%l] != bytestr[len(bytestr)-1] {
		return "", fmt.Errorf("密码或者校验位错误")
	}
	if base[((n-int(bytestr[0]))+r1)%l] != bytestr[0] {
		return "", fmt.Errorf("校验位错误")
	}
	var t1 = []byte(t[r1%l:] + t[:r1%l])
	for i := 0; i < r2%l; i++ {
		t1[i], t1[l-1-i] = t1[l-1-i], t1[i]
	}
	var mapping = map[byte]byte{}
	for i := 0; i < l; i++ {
		mapping[t1[i]] = base[i]
	}
	var e = []byte{}
	for _, char := range bytestr[1 : len(bytestr)-1] {
		v, ok := mapping[char]
		if !ok {
			return "", fmt.Errorf("字符集不匹配")
		}
		e = append(e, v)
	}
	return string(e), nil
}

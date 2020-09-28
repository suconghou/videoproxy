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
	Log = log.New(os.Stdout, "", 0)
)

// JSONPut resp json
func JSONPut(w http.ResponseWriter, v interface{}, status int, age int) (int, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
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

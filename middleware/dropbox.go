package middleware

import (
	"dropboxshare/util"
	"io"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/suconghou/go-dropbox"
)

var (
	client         *dropbox.Client
	publicDir      = "/Public"
	forwardHeaders = [...]string{"Range"}
	exposeHeaders  = [...]string{"Accept-Ranges", "Content-Range", "Connection", "Content-Length", "Content-Encoding", "Date", "Expires"}
	mimeTypes      = map[string]string{
		".mp4":  "video/mp4",
		".jpg":  "image/jpg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".json": "application/json",
		".js":   "text/javascript",
		".css":  "text/css",
		".html": "text/html",
	}
)

type fileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isdir"`
	Path    string `json:"path"`
	ModTime int64  `json:"mtime"`
}

type fileInfoList struct {
	List  []fileInfo `json:"list"`
	Total int        `json:"total"`
}

func init() {
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	if token == "" {
		token = ""
	}
	client = dropbox.New(token)
}

// Serve dropbox request
func Serve(w http.ResponseWriter, r *http.Request, uri string) error {
	var filePath = path.Join(publicDir, uri)
	file, err := client.Stat(filePath)
	if err != nil {
		return err
	}
	if file.IsDir() {
		fileList, err := client.List(filePath)
		if err != nil {
			return err
		}
		infoList := make([]fileInfo, 0)
		for _, item := range fileList {
			info := fileInfo{item.Name(), item.Size(), item.IsDir(), path.Join(uri, item.Name()), item.ModTime().Unix()}
			infoList = append(infoList, info)
		}
		res := fileInfoList{Total: len(fileList), List: infoList}
		util.JSONPut(w, res)
		return nil
	}
	resp, err := client.GetStream(filePath, func(h http.Header) {
		for _, item := range forwardHeaders {
			v := r.Header.Get(item)
			if v != "" {
				h.Set(item, v)
			}
		}
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	h := w.Header()
	for _, item := range exposeHeaders {
		v := resp.Header.Get(item)
		if v != "" {
			h.Set(item, v)
		}
	}
	ext := path.Ext(filePath)
	mtype := mime.TypeByExtension(ext)
	if mtype == "" {
		mtype = mimeTypes[ext]
	}
	if mtype != "" {
		h.Set("Content-Type", mtype)
	}
	w.WriteHeader(resp.StatusCode)
	if r.Method == "GET" {
		n, err := io.Copy(w, resp.Body)
		util.Logger.Printf("transfered %d bytes %v", n, err)
	}
	return nil
}

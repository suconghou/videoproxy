package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/suconghou/videoproxy/route"
	"github.com/suconghou/videoproxy/util"
)

func main() {
	var (
		port = flag.Int("p", 6060, "listen port")
		host = flag.String("h", "", "bind address")
	)
	flag.Parse()
	util.Log.Fatal(serve(*host, *port))
}

func serve(host string, port int) error {
	http.HandleFunc("/", routeMatch)
	util.Log.Printf("Starting up on port %d", port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	for _, p := range route.Route {
		if p.Reg.MatchString(r.URL.Path) {
			if err := p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path)); err != nil {
				util.Log.Print(err)
			}
			return
		}
	}
	fallback(w, r)
}

func fallback(w http.ResponseWriter, r *http.Request) {
	const index = "index.html"
	files := []string{index}
	if r.URL.Path != "/" {
		files = []string{r.URL.Path, path.Join(r.URL.Path, index)}
	}
	if !tryFiles(files, w, r) {
		if !tryFiles([]string{index}, w, r) {
			http.NotFound(w, r)
		}
	}
}

func tryFiles(files []string, w http.ResponseWriter, r *http.Request) bool {
	for _, file := range files {
		realpath := filepath.Join("./public", file)
		if f, err := os.Stat(realpath); err == nil {
			if f.Mode().IsRegular() {
				http.ServeFile(w, r, realpath)
				return true
			}
		}
	}
	return false
}

package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		var i map[string]interface{}
		d.Decode(&i)
		for k, v := range i {
			log.Printf("%s: %#v", k, v)
		}
		fmt.Fprintf(w, "Hello %q\n", html.EscapeString(r.URL.Path))
	case http.MethodGet:
		dir := http.Dir(os.Getenv("HOME"))
		path := html.EscapeString(r.URL.Path)
		f, err := dir.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "File not found: %s\n", path)
			return
		}
		stat, err := f.Stat()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			files, err := f.Readdir(0)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "<table>")
			fmt.Fprintf(w, "<tr><td><a href=\"..\">..</a></td></tr>")
			sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
			for _, dirFile := range files {
				name := dirFile.Name()
				if strings.HasPrefix(name, ".") && r.URL.Query().Get("showHidden") != "true" {
					continue
				}
				fmt.Fprintf(w, "<tr>")
				fmt.Fprintf(w, "<td>%6d</td>", dirFile.Size())
				if dirFile.IsDir() {
					name = fmt.Sprintf("%s/", name)
				}
				url := url.URL{Path: name}
				if path == "/" {
					fmt.Fprintf(w, "<td><a href=\"%s\">%s</a></td>", url.String(), name)
				} else {
					fmt.Fprintf(w, "<td><a href=\"%s%s\">%s</a></td>", path, url.String(), name)
				}
				fmt.Fprintf(w, "<td>%s</td>", dirFile.ModTime().Format("01/02/2006 15:04:05 MST"))
				fmt.Fprintf(w, "</tr>")
			}
			fmt.Fprintf(w, "</table>")
		} else {
			size := stat.Size()
			buf, err := io.ReadAll(f)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if n, err := w.Write(buf); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else if int64(n) != size {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	//case http.MethodDelete:
	//path := html.EscapeString(r.URL.Path)
	//err := os.Remove(strings.TrimPrefix(path, "/"))
	//if err != nil {
	//w.WriteHeader(http.StatusNotFound)
	//fmt.Fprintf(w, "File not found: %s\n", path)
	//return
	//}
	//fmt.Fprintf(w, "Deleted file: %s\n", path)
	default:
		fmt.Fprintf(w, "No\n")
	}
}

func main() {
	http.HandleFunc("/", handle)
	//dir := http.Dir(os.Getenv("HOME"))
	//fs := http.FileServer(dir)
	//http.Handle("/", fs)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
	select {}
}

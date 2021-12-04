package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		var i map[string]interface{}
		d.Decode(&i)
		for k, v := range i {
			switch vv := v.(type) {
			case string:
				fmt.Fprintf(w, "%s, %s\n", k, vv)
			case float64:
				fmt.Fprintf(w, "%s: %.0f\n", k, vv)
			}
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
			for _, dirFile := range files {
				fmt.Fprintf(w, "<tr>")
				fmt.Fprintf(w, "<td>%6d</td>", dirFile.Size())
				name := dirFile.Name()
				if dirFile.IsDir() {
					name = fmt.Sprintf("%s/", name)
				}
				if path == "/" {
					fmt.Fprintf(w, "<td><a href=\"%s\">%-30s</a></td>", name, name)
				} else {
					fmt.Fprintf(w, "<td><a href=\"%s%s\">%-30s</a></td>", path, name, name)
				}
				fmt.Fprintf(w, "<td>%s</td>", dirFile.ModTime().String())
				fmt.Fprintf(w, "</tr>")
			}
			fmt.Fprintf(w, "</table>")
			return
		}
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
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
	select {}
}

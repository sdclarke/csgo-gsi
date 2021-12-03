package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		d := json.NewDecoder(r.Body)
		var i map[string]interface{}
		d.Decode(&i)
		for _, v := range i {
			switch vv := v.(type) {
			case string:
				fmt.Fprintf(w, "Hello, %s\n", vv)
			case float64:
				fmt.Fprintf(w, "Age: %.0f\n", vv)
			}
		}
		fmt.Fprintf(w, "Hello %q\n", html.EscapeString(r.URL.Path))
	case http.MethodGet:
		var dir http.Dir = "/home/scottclarke"
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
		fmt.Fprintf(w, "What the actual fuck\n")
	}
}

func main() {
	http.HandleFunc("/", handle)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
	select {}
}

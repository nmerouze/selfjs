package main

import (
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/nmerouze/selfjs"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	bundle, _ := ioutil.ReadFile("./build/bundle.js")
	http.Handle("/", selfjs.New(runtime.NumCPU(), string(bundle)))
	http.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	http.Handle("/universal.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./build/bundle.js")
	}))

	http.ListenAndServe(":8080", nil)
}

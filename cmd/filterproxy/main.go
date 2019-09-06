package main

import (
	"log"
	"net/http"

	"github.com/wader/filtertransport/proxy"
)

func main() {
	http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		proxy.DefaultHandler.ServeHTTP(rw, r)
	}))
}

package main

import (
	"log"
	"net"
	"net/http"

	"github.com/wader/filtertransport"
)

func main() {
	proxyListener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("proxy listen err=%v", err)
	}
	log.Printf("listen on %s", proxyListener.Addr())
	http.Serve(proxyListener, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		filtertransport.DefaultHandler.ServeHTTP(rw, r)
	}))
}

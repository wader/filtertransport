package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wader/filtertransport"
)

func main() {
	c := http.Client{Transport: filtertransport.DefaultTransport}
	if r, err := c.Get(os.Args[1]); err != nil {
		log.Print(err)
	} else {
		defer r.Body.Close()
		log.Print(r.Status)
		for k, v := range r.Header {
			log.Printf("%s: %s", k, v)
		}
		io.Copy(os.Stdout, r.Body)
	}
}

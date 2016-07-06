package filtertransport

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

// Handler proxy handler filtering requests
type Handler struct {
	reverseProxy httputil.ReverseProxy
	filter       FilterTCPAddrFn
}

func copyAndClose(dst io.WriteCloser, src io.Reader) {
	io.Copy(dst, src)
	dst.Close()
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		// not connect, use ReverseProxy
		h.reverseProxy.ServeHTTP(rw, r)
		return
	}

	// connect, probably TLS, make TCP connection and tunnel traffic

	conn, err := DialTCP(r.Host, h.filter)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	hijack, ok := rw.(http.Hijacker)
	if !ok {
		// TODO: log hijack not supported?
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	hijackConn, _, err := hijack.Hijack()
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		hijackConn.Close()
		return
	}

	hijackConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	go copyAndClose(conn, hijackConn)
	go copyAndClose(hijackConn, conn)
}

// NewHandler new proxy handler using filter function to filter requests
func NewHandler(filter FilterTCPAddrFn) http.Handler {
	return &Handler{
		reverseProxy: httputil.ReverseProxy{
			ErrorLog: log.New(ioutil.Discard, "", 0), // prints to stderr if nil
			Director: func(r *http.Request) {},
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return DialTCP(addr, filter)
				},
			},
		},
		filter: filter,
	}
}

// DefaultHandler proxy that filters private addresses
var DefaultHandler = NewHandler(FilterPrivate)

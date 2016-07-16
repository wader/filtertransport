package filtertransport

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

// Handler proxy handler filtering requests
type Handler struct {
	reverseProxy httputil.ReverseProxy
	// used becuse ReverseProxy.Transport is a http.RoundTripper interface
	// we could cast but feels a bit ugly
	transport *http.Transport
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

	conn, err := h.transport.Dial("tcp", r.Host)
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

// NewHandler new proxy handler using transport to make actual requests
// use a transport with a filter dialer to filter proxy requests
func NewHandler(transport *http.Transport) http.Handler {
	return &Handler{
		reverseProxy: httputil.ReverseProxy{
			ErrorLog:  log.New(ioutil.Discard, "", 0), // prints to stderr if nil
			Director:  func(r *http.Request) {},
			Transport: transport,
		},
		transport: transport,
	}
}

// DefaultHandler proxy that filters private addresses
var DefaultHandler = NewHandler(DefaultTransport)

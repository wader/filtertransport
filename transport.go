package filtertransport

import (
	"net"
	"net/http"
)

// DefaultTransport http.DefaultTransport that filters private addresses
var DefaultTransport = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		return DialTCP(addr, FilterPrivate)
	},
}

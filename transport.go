package filtertransport

import (
	"net"
	"net/http"
	"time"
)

// DefaultTransport http.DefaultTransport that filters using DefaultFilter
var DefaultTransport = &http.Transport{
	// does not include ProxyFromEnvironment, makes no sense for filter
	Dial: func(network, addr string) (net.Conn, error) {
		return FilterDial(network, addr, DefaultFilter, (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial)
	},
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

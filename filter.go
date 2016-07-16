/*
Package filtertransport implements filtering http transport and proxy handler
*/
package filtertransport

import (
	"fmt"
	"net"
)

// DialFn http.Transport dial function
type DialFn func(network string, address string) (net.Conn, error)

// FilterTCPAddrFn function deciding if to filter
type FilterTCPAddrFn func(addr net.TCPAddr) error

// FilterError TCP address filtered error
type FilterError struct {
	net.TCPAddr
}

func (e FilterError) Error() string {
	return fmt.Sprintf("%s is filtered", e.TCPAddr.String())
}

// MustParseCIDR parses string into net.IPNet
func MustParseCIDR(s string) net.IPNet {
	if _, ipnet, err := net.ParseCIDR(s); err != nil {
		panic(err)
	} else {
		return *ipnet
	}
}

// PrivateNetworks private net.IPNets
var PrivateNetworks = []net.IPNet{
	MustParseCIDR("169.254.0.0/16"),
	MustParseCIDR("172.16.0.0/12"),
	MustParseCIDR("192.168.0.0/16"),
	MustParseCIDR("10.0.0.0/8"),
	MustParseCIDR("127.0.0.0/8"),
	MustParseCIDR("::1/128"),
	MustParseCIDR("fc00::/7"),
}

// FilterDial http.Transport dial with filtering function
func FilterDial(network string, address string, filter FilterTCPAddrFn, dial DialFn) (net.Conn, error) {
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, fmt.Errorf("unsupported network %s", network)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	if err := filter(*tcpAddr); err != nil {
		return nil, err
	}

	// pass along resolved address to prevent DNS rebind
	return dial(network, tcpAddr.String())
}

func findIPNet(ipnets []net.IPNet, ip net.IP) bool {
	for _, ipnet := range ipnets {
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

// FilterPrivate filter function filtering PrivateNetworks
func FilterPrivate(addr net.TCPAddr) error {
	if findIPNet(PrivateNetworks, addr.IP) {
		return FilterError{addr}
	}
	return nil
}

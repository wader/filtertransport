// Package filtertransport implements filtering http transport and proxy handler
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

// DefaultFilteredNetworks net.IPNets that are loopback, private, link local, default unicast
var DefaultFilteredNetworks = []net.IPNet{
	MustParseCIDR("127.0.0.0/8"),    // loopback
	MustParseCIDR("0.0.0.0/32"),     // default unicast
	MustParseCIDR("169.254.0.0/16"), // link local
	MustParseCIDR("172.16.0.0/12"),  // private
	MustParseCIDR("192.168.0.0/16"), // private
	MustParseCIDR("10.0.0.0/8"),     // private
	MustParseCIDR("::1/128"),        // loopback
	MustParseCIDR("::/128"),         // default unicast
	MustParseCIDR("fc00::/7"),       // unique local addresses
	// IPv4 compatibility network (overlaps ::/128)
	// not really private but very rare
	MustParseCIDR("::/96"),
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

// FindIPNet true if any of the ipnets contains ip
func FindIPNet(ipnets []net.IPNet, ip net.IP) bool {
	for _, ipnet := range ipnets {
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

// DefaultFilter filters DefaultFilteredNetworks
func DefaultFilter(addr net.TCPAddr) error {
	if FindIPNet(DefaultFilteredNetworks, addr.IP) {
		return FilterError{addr}
	}
	return nil
}

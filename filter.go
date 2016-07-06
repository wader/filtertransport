/*
Package filtertransport implements filtering http transport and proxy handler
*/
package filtertransport

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

// FilterTCPAddrFn is function deciding if to filter
type FilterTCPAddrFn func(addr net.TCPAddr) error

// FilterError request filtered error
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

// DialTCP like net.DialTCP but with filtering function
func DialTCP(address string, filter FilterTCPAddrFn) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	if err := filter(*tcpAddr); err != nil {
		return nil, err
	}

	var af int
	var sa syscall.Sockaddr

	if ip4 := tcpAddr.IP.To4(); ip4 != nil {
		af = syscall.AF_INET
		sa4 := &syscall.SockaddrInet4{Port: tcpAddr.Port}
		copy(sa4.Addr[:], ip4)
		sa = sa4
	} else if ip16 := tcpAddr.IP.To16(); ip16 != nil {
		af = syscall.AF_INET6
		sa6 := &syscall.SockaddrInet6{Port: tcpAddr.Port}
		copy(sa6.Addr[:], ip16)
		sa = sa6
	} else {
		return nil, fmt.Errorf("unknown ip len %d (%s)", len(tcpAddr.IP), address)
	}

	fd, err := syscall.Socket(af, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}

	if err := syscall.Connect(fd, sa); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "")
	return net.FileConn(file)
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

/*
 * Author: 空闻
 * Date: 2021/11/4
 * Description: .
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
	"net"
)

// ip route get 193.233.7.82
func ipRoute(proto string, desIP net.IP, lport string) (net.Addr, error) {
	r, err := netlink.RouteGet(desIP)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, errors.New("empty")
	}

	srcIP := r[0].Src
	if srcIP.IsLinkLocalUnicast() || srcIP.IsLinkLocalMulticast() || srcIP.IsLoopback() {
		return nil, errors.New("address err: " + srcIP.String())
	}

	var addr net.Addr
	switch proto {
	case "udp":
		addr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(srcIP.String(), lport))
	case "tcp":
		addr, err = net.ResolveTCPAddr("tcp", net.JoinHostPort(srcIP.String(), lport))
	default:
		return nil, errors.New("protocol err: " + proto)
	}

	if err != nil {
		return nil, err
	}

	return addr, nil
}

func dumpRoue() ([]netlink.Route, error){
	return netlink.RouteList(nil, nl.FAMILY_V4)
}


func main() {
	var ipStr string
	var port string
	flag.StringVar(&ipStr, "ip", "8.8.8.8", "dest ip")
	flag.StringVar(&port, "src port", "80", "src port")
	flag.Parse()

	dst := net.ParseIP(ipStr)

	lAddr, e :=ipRoute("udp", dst, "80")
	if e == nil {
		fmt.Println("[src]", lAddr, " ->  [dest]", ipStr)
	} else {
		fmt.Println("ip route error: ", e)
	}


	fmt.Println(dumpRoue())
}

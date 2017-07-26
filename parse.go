package tcpping

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

const (
	separatorSign = ","
)

func ParseHosts(in string) (hosts []string) {
	items := strings.Split(in, separatorSign)
	for _, it := range items {
		hosts = append(hosts, parseHost(it)...)
	}
	return
}

func parseHost(in string) []string {
	ip, ipNet, err := net.ParseCIDR(in)
	if err == nil {
		ips := []string{}
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			ips = append(ips, ip.String())
		}
		return ips[1 : len(ips)-1]
	}
	ip = net.ParseIP(in)
	if ip != nil {
		return []string{ip.String()}
	}
	return []string{in}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ParsePorts(in []string) (ports []uint16, err error) {
	items := strings.Split(in, separatorSign)
	for _, it := range items {
		ports, err = parsePortAndAppend(ports, it)
		if err != nil {
			return
		}
	}
	return
}

func parsePortAndAppend(ports []uint16, in string) ([]uint16, error) {
	port, err := parsePort(in)
	if err != nil {
		return ports, err
	}
	return append(ports, port), nil
}

const (
	minPort = 0
	maxPort = 65535
)

var ErrPortOutOfRange = errors.New("Port out of range")

func parsePort(in string) (uint16, error) {
	n, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if n > maxPort || n < minPort {
		return 0, ErrPortOutOfRange
	}
	return uint16(n), nil
}

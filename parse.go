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
		return ips
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

func ParsePorts(in string) (ports []uint16, err error) {
	if in == "" {
		return
	}
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
	parsedPorts, err := parsePorts(in)
	if err != nil {
		return ports, err
	}
	return append(ports, parsedPorts...), nil
}

const (
	minPort        = 0
	maxPort        = 65535
	portsRangeSign = "-"
)

var (
	ErrPortOutOfRange    = errors.New("Port out of range")
	ErrInvalidPort       = errors.New("Invalid port")
	ErrInvalidPortsRange = errors.New("Invalid ports range")
)

func parsePorts(in string) ([]uint16, error) {
	parts := strings.Split(in, portsRangeSign)
	switch len(parts) {
	case 2:
		return parsePortsRange(parts)
	case 1:
		port, err := parseSinglePort(in)
		return []uint16{port}, err
	default:
		return nil, ErrInvalidPort
	}
}

func parsePortsRange(parts []string) ([]uint16, error) {
	points := []uint16{}
	for _, part := range parts {
		point, err := parseSinglePort(part)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	if points[0] > points[1] {
		return nil, ErrInvalidPortsRange
	}
	ports := []uint16{}
	for i := points[0]; i < points[1]+1; i++ {
		ports = append(ports, i)
	}
	return ports, nil
}

func parseSinglePort(in string) (uint16, error) {
	n, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if n > maxPort || n < minPort {
		return 0, ErrPortOutOfRange
	}
	return uint16(n), nil
}

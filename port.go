package tcpping

import (
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Port struct {
	Port   uint16
	Status PortStatus
}

type PortStatus int

const (
	PortStatusClosed PortStatus = iota
	PortStatusOpen
	PortStatusTimeout
)

func (e PortStatus) String() string {
	switch e {
	case PortStatusClosed:
		return "closed"
	case PortStatusOpen:
		return "open"
	case PortStatusTimeout:
		return "timeout"
	default:
		return ""
	}
}

const (
	portColumnLength = 6
)

func (p *Port) String() string {
	out := strconv.Itoa(int(p.Port))
	out += strings.Repeat(" ", portColumnLength-len(out))
	return out + p.Status.String()
}

const (
	dialTimeout = time.Second
	tcpNetwork  = "tcp"
)

func (p *Port) check(host string) error {
	conn, err := net.DialTimeout(tcpNetwork,
		net.JoinHostPort(host, strconv.Itoa(int(p.Port))), dialTimeout)

	switch e := err.(type) {
	case *net.OpError:
		if sysErr, ok := e.Err.(*os.SyscallError); ok && sysErr.Err == syscall.ECONNREFUSED {
			p.Status = PortStatusClosed
		} else if e.Timeout() {
			p.Status = PortStatusTimeout
		} else {
			return err
		}
	case nil:
		conn.Close()
		p.Status = PortStatusOpen
	default:
		return err
	}
	return nil
}

type PortChecker struct {
	hosts   []string
	ports   []uint16
	verbose bool
	sync.Mutex
}

func (t *PortChecker) SetHosts(hosts []string) *PortChecker {
	t.Lock()
	t.hosts = hosts
	t.Unlock()
	return t
}

func (t *PortChecker) SetPorts(ports []uint16) *PortChecker {
	t.Lock()
	t.ports = ports
	t.Unlock()
	return t
}

func (t *PortChecker) SetVerbose(v bool) *PortChecker {
	t.Lock()
	t.verbose = v
	t.Unlock()
	return t
}

type PortCheckResult struct {
	Port Port
	Err  error
}

func (p *PortCheckResult) String() string {
	if p.Err != nil {
		return p.Err.Error()
	}
	return p.Port.String()
}

const (
	defaultEndRangePort   = 1024
	defaultStartRangePort = 0
)

func (t *PortChecker) Run(closeChan chan struct{}) chan PortCheckResult {
	t.Lock()
	if len(t.ports) == 0 {
		for i := defaultStartRangePort; i < defaultEndRangePort; i++ {
			t.ports = append(t.ports, uint16(i))
		}
	}
	resultChan := make(chan PortCheckResult, 1)
	go t.run(resultChan, closeChan)
	t.Unlock()
	return resultChan
}

func (t *PortChecker) run(resultChan chan PortCheckResult, closeChan chan struct{}) {
	t.Lock()
	defer close(resultChan)
	defer t.Unlock()

	for _, host := range t.hosts {
		for _, port := range t.ports {
			select {
			case <-closeChan:
				return
			default:
				it := check(host, port)
				if it.Port.Status == PortStatusTimeout && !t.verbose {
					continue
				}
				resultChan <- it
			}
		}
	}
}

func check(host string, port uint16) (res PortCheckResult) {
	res.Port.Port = port
	res.Err = res.Port.check(host)
	return
}

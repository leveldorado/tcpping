
## Overview

tcpping is a command line tool and lib for network exploration.

Install lib with:

    $ go get github.com/leveldorado/tcpping


Lib usage:

```go

func main() {
    portChecker := tcpping.PortChecker{}
  	portChecker.SetPorts([]uint16{23,80}).SetHosts([]string{"127.0.0.1"}).SetVerbose(true)
  	resultChan := portChecker.Run(nil)
  	for it := range resultChan {
  		fmt.Println(it.String())
  	}
}
```

Install tool with:

    $ go install github.com/leveldorado/tcpping/cmd/tcpping

Tool usage:

```
$ tcpping 127.0.0.1,google.com,134.5.6.3/30 -p 80,400-405 -v
```


### Flags and args

first arg target could be ip, domain, CIDR.

-v --verbose shows in list timeout ports

-p --ports  could be single port, range of ports e.g. -p 40-45
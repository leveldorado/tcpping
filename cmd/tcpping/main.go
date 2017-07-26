package main

import (
	"fmt"
	"github.com/leveldorado/tcpping"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	opt := parseOptions()
	ports, err := tcpping.ParsePorts(opt.ports)
	if err != nil {
		log.Fatal(err)
	}
	t := tcpping.PortChecker{}
	t.SetHosts(tcpping.ParseHosts(opt.targets)).SetVerbose(opt.verbose).SetPorts(ports)

	closeChan := make(chan struct{}, 1)
	resultChan := t.Run(closeChan)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				return
			}
			fmt.Println(result.String())
		case _ = <-c:
			closeChan <- struct{}{}
		}
	}
}

func parseOptions() (opt options) {
	kingpin.Arg("target", "target or targets to explore").Required().StringVar(&opt.targets)
	kingpin.Flag("ports", "port or ports to explore").Short('p').StringVar(&opt.ports)
	kingpin.Flag("verbose", "if set include timeout ports").Short('v').BoolVar(&opt.verbose)
	kingpin.Parse()
	return
}

type options struct {
	targets string
	ports   string
	verbose bool
}

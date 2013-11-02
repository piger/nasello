// Nasello is a DNS proxy server
package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/piger/nasello"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	configFile = flag.String("config", "nasello.json", "Configuration file")
)

func serve(net string, address string) {
	err := dns.ListenAndServe(address, net, nil)
	if err != nil {
		log.Fatalf("Failed to setup the "+net+" server: %s\n", err.Error())
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	flag.Parse()

	configuration := nasello.ReadConfig(*configFile)

	for _, filter := range configuration.Filters {
		// Ensure that the pattern is a FQDN name
		pattern := dns.Fqdn(filter.Pattern)

		log.Printf("Proxing %s on %v\n", pattern, filter.Addresses)
		dns.HandleFunc(pattern, nasello.ServerHandler(filter.Addresses))
	}
	go serve("tcp", ":8053")
	go serve("udp", ":8053")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
forever:
	for {
		select {
		case s := <-sig:
			log.Printf("Signal (%d) received, stopping\n", s)
			break forever
		}
	}
}

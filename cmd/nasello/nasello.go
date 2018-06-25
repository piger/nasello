// Nasello is a DNS proxy server
package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/piger/nasello"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	configFile = flag.String("config", "nasello.json", "Configuration file")
	listenAddr = flag.String("listen", "localhost:8053", "Local bind address")
)

func serve(net string, address string) {
	err := dns.ListenAndServe(address, net, nil)
	if err != nil {
		log.Fatalf("Failed to setup the "+net+" server: %s\n", err.Error())
	}
}

func main() {
	flag.Parse()

	configuration := nasello.ReadConfig(*configFile)
	for _, filter := range configuration.Filters {
		// Ensure that each pattern is a FQDN name
		pattern := dns.Fqdn(filter.Pattern)

		log.Printf("Proxing %s on %v(%s)\n", pattern, strings.Join(filter.Addresses, ", "), filter.Protocol)
		dns.HandleFunc(pattern, nasello.ServerHandler(filter.Addresses, filter.Protocol))
	}

	go serve("tcp", *listenAddr)
	go serve("udp", *listenAddr)

	log.Printf("Started DNS server on: %s\n", *listenAddr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-sig
	log.Printf("Signal (%s) received, stopping\n", recvSig.String())
}

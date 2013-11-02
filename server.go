// Nasello is a DNS proxy server.
//
// It can be used to route DNS queries to different remote servers based on
// pattern matching on the requested name.
//
// See `config.go` for details about the configuration file format.
//
// Code is inspired by go-dns examples like:
// https://github.com/miekg/exdns/blob/master/q/q.go
package nasello

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

type handler func(dns.ResponseWriter, *dns.Msg)

// Returns an anonymous function configured to resolve DNS
// queries with a specific set of remote servers.
func ServerHandler(addresses []string) handler {

	return func(w dns.ResponseWriter, req *dns.Msg) {
		nameserver := addresses[0]

		if !strings.Contains(nameserver, ":") {
			nameserver = net.JoinHostPort(nameserver, "53")
		}

		log.Printf("Incoming request: %s %s %v - using %s\n",
			dns.ClassToString[req.Question[0].Qclass],
			dns.TypeToString[req.Question[0].Qtype],
			req.Question[0].Name, nameserver)

		c := new(dns.Client)
		c.Net = "udp"

		resp, rtt, err := c.Exchange(req, nameserver)

		for {

		Redo:
			if err != nil {
				log.Printf(";; ERROR: %s\n", err.Error())
				continue
			}

			if req.Id != resp.Id {
				log.Printf("Id mismatch: %v != %v\n", req.Id, resp.Id)
				return
			}

			if resp.MsgHdr.Truncated {
				log.Printf("Truncated message, retrying TCP")
				c.Net = "tcp"
				resp, rtt, err = c.Exchange(req, nameserver)
				goto Redo
			}

			log.Printf("Query time: %.3d µs, server: %s(%s), size: %d bytes\n", rtt/1e3, nameserver, c.Net, resp.Len())
			w.WriteMsg(resp)
			return
		}
	}
}

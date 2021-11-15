// Package nasello provides a simple DNS proxy server.
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
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/miekg/dns"
)

// Handler is the handler function that will serve DNS requests.
type Handler func(dns.ResponseWriter, *dns.Msg)

// ServerHandler Returns an anonymous function configured to resolve DNS
// queries with a specific set of remote servers.
func ServerHandler(addresses []string, protocol string) Handler {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// This is the actual handler
	return func(w dns.ResponseWriter, req *dns.Msg) {
		nameserver := addresses[randGen.Intn(len(addresses))]

		if protocol == "" {
			switch t := w.RemoteAddr().(type) {
			default:
				log.Printf("ERROR: Unsupported protocol %T\n", t)
				return
			case *net.UDPAddr:
				protocol = "udp"
			case *net.TCPAddr:
				protocol = "tcp"
			}
		}

		for _, q := range req.Question {
			log.Printf("Incoming request #%v: %s %s %v - using %s(%s)\n",
				req.Id,
				dns.ClassToString[q.Qclass],
				dns.TypeToString[q.Qtype],
				q.Name, nameserver, protocol)
		}

		c := new(dns.Client)
		c.Net = protocol
		resp, rtt, err := c.Exchange(req, nameserver)

	Redo:
		switch {
		case err != nil:
			log.Printf("ERROR: %s\n", err.Error())
			sendFailure(w, req)
			return
		case req.Id != resp.Id:
			log.Printf("ERROR: Id mismatch: %v != %v\n", req.Id, resp.Id)
			sendFailure(w, req)
			return
		case resp.MsgHdr.Truncated && protocol != "tcp":
			log.Printf("WARNING: Truncated answer for request %v, retrying TCP\n", req.Id)
			c.Net = "tcp"
			resp, rtt, err = c.Exchange(req, nameserver)
			goto Redo
		}

		log.Printf("Request #%v: %.3d µs, server: %s(%s), size: %d bytes\n", resp.Id, rtt/1e3, nameserver, c.Net, resp.Len())
		if err := w.WriteMsg(resp); err != nil {
			log.Printf("ERROR: write failed: %s", err)
		}
	} // end of handler
}

func sendFailure(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetRcode(r, dns.RcodeServerFailure)
	if err := w.WriteMsg(msg); err != nil {
		log.Printf("ERROR: write failed in sendFailure: %s", err)
	}
}

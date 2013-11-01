package nasello

import (
	"log"
	"net"
	"strings"
	"github.com/miekg/dns"
)

type handler func(dns.ResponseWriter, *dns.Msg)

func ServerHandler(addresses []string) handler {

	return func (w dns.ResponseWriter, r *dns.Msg) {
		nameserver := addresses[0]

		if (!strings.Contains(nameserver, ":")) {
			nameserver = net.JoinHostPort(nameserver, "53")
		}

		log.Printf("Incoming request: %s %s %v - using %s\n",
			dns.ClassToString[r.Question[0].Qclass],
			dns.TypeToString[r.Question[0].Qtype],
			r.Question[0].Name, nameserver)

		c := new(dns.Client)
		c.Net = "udp"

		resp, rtt, err := c.Exchange(r, nameserver)

		for {

		Redo:
			if err != nil {
				log.Printf(";; ERROR: %s\n", err.Error())
				continue
			}

			if r.MsgHdr.Truncated {
				log.Printf("Truncated message, retrying TCP")
				c.Net = "tcp"
				resp, rtt, err = c.Exchange(r, nameserver)
				goto Redo
			} 
			
			log.Printf("Answered in %v\n", rtt)
			w.WriteMsg(resp)
			return
		}
	}
}

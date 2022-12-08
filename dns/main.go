package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"redes.cucho.": "192.168.0.2",
}

func parseQuery(msg *dns.Msg) {
	for _, question := range msg.Question {
		if question.Qtype == dns.TypeA {
			log.Printf("Query for %s\n", question.Name)

			ip := records[question.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", question.Name, ip))
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
				}
			}
		}
	}
}

func handleDNSRequest(writer dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	if r.Opcode == dns.OpcodeQuery {
		parseQuery(msg)
	}

	if err := writer.WriteMsg(msg); err != nil {
		log.Printf("Failed to write response: %s\n", err.Error())
	}
}

func main() {
	// attach request handler func
	dns.HandleFunc("cucho.", handleDNSRequest)

	// start server
	port := 5353
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}

	log.Printf("Starting at %d\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to shutdown server: %s\n ", err.Error())
		}
	}()
}
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var records = map[string]string{}

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

func parseRecordsFile(filename string) error {
	file, err := os.Open(filename) // #nosec
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %s\n", err.Error())
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			records[parts[0]] = parts[1]
		}
	}

	return nil
}

func main() {
	// attach request handler func
	dns.HandleFunc("cucho.", handleDNSRequest)

	port := "53"
	if len(os.Args) > 1 && os.Args[1] != "" {
		port = os.Args[1]
	}

	recordsFile := "records.txt"
	if len(os.Args) > 2 && os.Args[2] != "" {
		recordsFile = os.Args[2]
	}

	err := parseRecordsFile(recordsFile)
	if err != nil {
		log.Fatalf("Failed to parse records file: %s\n", err.Error())
	}

	server := &dns.Server{Addr: ":" + port, Net: "udp"}

	log.Printf("Starting at %s\n", port)

	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}

	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to shutdown server: %s\n", err.Error())
		}
	}()
}

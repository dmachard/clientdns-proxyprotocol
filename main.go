package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/miekg/dns"
	"github.com/pires/go-proxyproto"
)

// Function to build a Proxy Protocol v2 header
func buildProxyHeader(srcAddr net.Addr, dstAddr net.Addr, kv map[uint8][]byte) ([]byte, error) {
	protocol := proxyproto.UDPv4

	header := proxyproto.Header{
		Version:           2,
		Command:           proxyproto.PROXY,
		TransportProtocol: protocol,
		SourceAddr:        srcAddr,
		DestinationAddr:   dstAddr,
	}

	for keyInt, value := range kv {
		keyHex := proxyproto.PP2Type(keyInt & 0xFF)
		tlvs := []proxyproto.TLV{
			{
				Type:  keyHex,
				Value: value,
			},
		}
		err := header.SetTLVs(tlvs)
		if err != nil {
			return nil, fmt.Errorf("Error to set TLVs: %w", err)
		}
	}

	// Serialize the Proxy Protocol header into a byte slice
	headerBytes, err := header.Format()
	if err != nil {
		return nil, fmt.Errorf("failed to format Proxy Protocol header: %w", err)
	}

	// Combine the Proxy Protocol header and TLVs into a single byte slice
	return headerBytes, nil
}

// Function to send a DNS query with Proxy Protocol
func sendDNSQuery(server string, port string, domain string, qtype uint16, kv map[uint8][]byte) {
	// Build the DNS request
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)

	// Resolve the DNS server address
	protocol := "udp"
	serverAddress := net.JoinHostPort(server, port)
	dstAddr, err := net.ResolveUDPAddr(protocol, serverAddress)
	if err != nil {
		log.Fatalf("Error resolving DNS server: %v", err)
	}

	// Connect to the DNS server
	conn, err := net.Dial(protocol, serverAddress)
	if err != nil {
		log.Fatalf("Error connecting to DNS server: %v", err)
	}
	defer conn.Close()

	// Get the source address
	srcAddr := conn.LocalAddr()

	// Build the Proxy Protocol header
	proxyHeader, err := buildProxyHeader(srcAddr, dstAddr, kv)
	if err != nil {
		log.Fatalf("Error building Proxy Protocol header: %v", err)
	}

	// Serialize the DNS query
	dnsBytes, err := msg.Pack()
	if err != nil {
		log.Fatalf("Error serializing DNS query: %v", err)
	}

	// Concatenate Proxy Protocol header + DNS query
	fullPacket := append(proxyHeader, dnsBytes...)

	// Send the DNS query
	_, err = conn.Write(fullPacket)
	if err != nil {
		log.Fatalf("Error sending query: %v", err)
	}

	// Read the response
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Decode the response
	resp := new(dns.Msg)
	err = resp.Unpack(buffer[:n])
	if err != nil {
		log.Fatalf("Error unpacking DNS response: %v", err)
	}

	// Print the response
	fmt.Println(resp)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <dns_server> <port> <domain> <type(A,AAAA,MX)> [kv_key=kv_value]\n", os.Args[0])
		os.Exit(1)
	}

	server := os.Args[1]
	port := os.Args[2]
	domain := os.Args[3]
	qtype := dns.TypeA
	kv := make(map[uint8][]byte)

	// Determine query type
	switch os.Args[4] {
	case "AAAA":
		qtype = dns.TypeAAAA
	case "MX":
		qtype = dns.TypeMX
	}

	// Parse Key-Value arguments (TLVs)
	for _, arg := range os.Args[5:] {
		var keyInt int
		var value string
		_, err := fmt.Sscanf(arg, "%d=%s", &keyInt, &value)
		if err == nil {
			kv[uint8(keyInt)] = []byte(value)
		}
	}

	// Send the DNS query
	sendDNSQuery(server, port, domain, qtype, kv)
}

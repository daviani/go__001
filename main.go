package main

import (
	"fmt"

	"github.com/daviani/go__001/internal/scanner"
)

func main() {

	dns := scanner.DNSScanner{}
	ssl := scanner.SSLScanner{}
	fmt.Printf(dns.Scan("daviani.dev"))
	fmt.Printf(ssl.Scan("daviani.dev------"))

	runScann(ssl, "daviani.dev")
	runScann(dns, "daviani.dev")
}

func runScann(s scanner.Scanner, domain string) {
	fmt.Printf(s.Scan(domain))
}

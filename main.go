package main

import (
	"fmt"

	"github.com/daviani/go__001/internal/scanner"
)

func main() {

	dns := scanner.DNSScanner{}
	ssl := scanner.SSLScanner{}
	scanners := []scanner.Scanner{dns, ssl}
	results := map[string]string{}

	for _, value := range scanners {
		result := value.Scan("daviani.dev")
		key := value.Name() + "_daviani.dev"
		results[key] = result
	}
	fmt.Println(results)
}

func runScann(s scanner.Scanner, domain string) {
	fmt.Println(s.Scan(domain))
}

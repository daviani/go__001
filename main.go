package main

import (
	"github.com/daviani/go__001/internal/api"
	"github.com/daviani/go__001/internal/scanner"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}

// @title           GoSentry — Security Audit API
// @version         1.0
// @description     API d'audit de surface d'attaque externe. Analyse DNS, certificats SSL/TLS, headers de sécurité, sous-domaines et fichiers sensibles exposés.
// @host            localhost:8082
// @BasePath        /
func main() {
	// Initialisation des scanners (structs vides qui implémentent l'interface Scanner)
	dns := scanner.DNSScanner{}
	ssl := scanner.SSLScanner{}
	header := scanner.HeaderScanner{}
	subdomain := scanner.SubdomainScanner{}
	sensitive := scanner.SensitiveScanner{}

	// Slice contenant tous les scanners - on peut en ajouter autant qu'on veut
	scanners := []scanner.Scanner{dns, ssl, header, subdomain, sensitive}

	server := api.Server{Port: 8082, Scanners: scanners}
	server.Start()
}

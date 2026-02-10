package api

import "github.com/daviani/go__001/internal/scanner"

// Server contient la configuration du serveur HTTP et la liste des scanners disponibles
type Server struct {
	Port     int               // Port d'écoute (ex: 8082)
	Scanners []scanner.Scanner // Slice des scanners — utilisée par handleAll pour les goroutines
}

// HealthResult — réponse JSON pour GET /health
type HealthResult struct {
	Status string `json:"status"`
}

// ScanResult — réponse JSON pour les routes /scan/*
type ScanResult struct {
	Scanner string `json:"scanner"` // Nom du scanner (dns, ssl, header...)
	Domain  string `json:"domain"`  // Domaine scanné
	Result  string `json:"result"`  // Résultat du scan (texte brut)
}

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/daviani/go__001/internal/scanner"
)

// Server contient la configuration du serveur HTTP
// Port : port d'écoute (ex: 8082)
// Scanners : slice d'interfaces Scanner disponibles pour les routes /scan/*
type Server struct {
	Port     int
	Scanners []scanner.Scanner
}

// HealthResult — structure de réponse JSON pour /health
// Sérialisée en JSON via json.NewEncoder : {"status": "ok"}
type HealthResult struct {
	Status string `json:"status"`
}

// ScanResult — structure de réponse JSON pour les routes /scan/*
// Les tags `json:"..."` définissent les noms des clés dans le JSON de sortie
// Équivalent JS : { scanner: "dns", domain: "daviani.dev", result: "..." }
type ScanResult struct {
	Scanner string `json:"scanner"`
	Domain  string `json:"domain"`
	Result  string `json:"result"`
}

// makeScanHandler retourne un handler HTTP pour un scanner donné (closure)
// Évite la duplication : le même code gère les 5 routes /scan/*
// name et s sont "capturés" par la closure — accessibles à chaque requête
// Équivalent JS : const makeHandler = (name, scanner) => (req, res) => { ... }
func makeScanHandler(name string, s scanner.Scanner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Query().Get("domain") extrait le query param "domain" de l'URL
		// Équivalent Express : req.query.domain
		domain := r.URL.Query().Get("domain")
		// Validation : si le domain est vide, retourne une erreur 400 (Bad Request)
		// 400 = erreur client ("tu as mal appelé l'API")
		// 500 = erreur serveur ("le serveur a planté")
		if domain == "" {
			http.Error(w, "paramètre 'domain' requis", http.StatusBadRequest)
			return
		}

		// Lance le scan DNS — peut échouer si le domaine est invalide
		result, err := s.Scan(domain)
		if err != nil {
			// http.Error envoie un message d'erreur + code HTTP 500 au client
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		// Encode le résultat dans le struct ScanResult et l'envoie en JSON
		err = json.NewEncoder(w).Encode(ScanResult{
			Scanner: name,
			Domain:  domain,
			Result:  result,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// Start enregistre les routes HTTP et démarre le serveur
// w http.ResponseWriter = équivalent de res en Express (on écrit la réponse dedans)
// r *http.Request = équivalent de req en Express (la requête entrante)
// Les routes sont enregistrées AVANT ListenAndServe qui est bloquant
func (s *Server) Start() {

	// Route /health — status du serveur
	// Retourne {"status": "ok"} pour vérifier que le serveur tourne
	http.HandleFunc(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			// Indique au client que la réponse est du JSON (pas du HTML ou du texte)
			w.Header().Set("Content-type", "application/json")
			// json.NewEncoder(w).Encode() sérialise le struct en JSON
			// et l'écrit directement dans le ResponseWriter
			// Équivalent Express : res.json({ status: "ok" })
			err := json.NewEncoder(w).Encode(HealthResult{Status: "ok"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		},
	)

	http.HandleFunc("/scan/dns", makeScanHandler("dns", scanner.DNSScanner{}))

	http.HandleFunc("/scan/ssl", makeScanHandler("ssl", scanner.SSLScanner{}))

	http.HandleFunc("/scan/header", makeScanHandler("header", scanner.HeaderScanner{}))

	http.HandleFunc("/scan/sensitive", makeScanHandler("sensitive", scanner.SensitiveScanner{}))

	http.HandleFunc("/scan/subdomain", makeScanHandler("subdomain", scanner.SubdomainScanner{}))

	// Démarrage du serveur — ListenAndServe est bloquant
	// Le programme reste ici et écoute les connexions entrantes
	fmt.Println("Serveur démarré sur le port :", s.Port)

	// fmt.Sprintf(":%d", s.Port) convertit l'int en string formatée (ex: 8082 → ":8082")
	// nil = utilise le routeur par défaut (celui où on a enregistré les routes avec HandleFunc)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.Port),
		nil,
	)

	// Si ListenAndServe retourne, c'est qu'il y a eu une erreur (ex: port déjà pris)
	// log.Fatal affiche l'erreur ET arrête le programme (exit code 1)
	if err != nil {
		log.Fatal(err)
	}
}

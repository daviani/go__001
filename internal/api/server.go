package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/daviani/go__001/docs"           // Blank import — enregistre la spec Swagger au démarrage via init()
	"github.com/daviani/go__001/internal/scanner" // Package contenant l'interface Scanner et ses implémentations
	httpSwagger "github.com/swaggo/http-swagger"  // Middleware servant l'interface Swagger UI
)

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

// makeScanHandler — closure qui retourne un handler HTTP pour un scanner donné
// Évite la duplication de code : le même pattern gère les 5 routes /scan/*
// name et s sont "capturés" par la closure et accessibles à chaque requête
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

		// Lance le scan — peut échouer si le domaine est invalide ou injoignable
		result, err := s.Scan(domain)
		if err != nil {
			log.Println(err)
			http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Encode le résultat dans le struct ScanResult et l'envoie en JSON
		err = json.NewEncoder(w).Encode(ScanResult{
			Scanner: name,
			Domain:  domain,
			Result:  result,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
			return
		}

	}
}

// @Summary     Status du serveur
// @Description Vérifie que le serveur est en ligne
// @Tags        health
// @Produce     json
// @Success     200 {object} HealthResult
// @Failure     500 {string} string "erreur serveur"
// @Router      /health [get]
func handleHealth() http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request) {
		// Indique au client que la réponse est du JSON (pas du HTML ou du texte)
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode() sérialise le struct en JSON
		// et l'écrit directement dans le ResponseWriter
		// Équivalent Express : res.json({ status: "ok" })
		err := json.NewEncoder(w).Encode(HealthResult{Status: "ok"})
		if err != nil {
			log.Println(err)
			http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
			return
		}
	}
}

// @Summary     Scan DNS
// @Description Analyse les records DNS du domaine (A, AAAA, MX, NS, TXT)
// @Tags        scanner
// @Produce     json
// @Param       domain query string true "Domaine à scanner"
// @Success     200 {object} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/dns [get]
func handleDNS() http.HandlerFunc {
	return makeScanHandler("dns", scanner.DNSScanner{})
}

// @Summary     Scan SSL/TLS
// @Description Analyse le certificat TLS du domaine (émetteur, expiration, validité)
// @Tags        scanner
// @Produce     json
// @Param       domain query string true "Domaine à scanner"
// @Success     200 {object} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/ssl [get]
func handleSSL() http.HandlerFunc {
	return makeScanHandler("ssl", scanner.SSLScanner{})
}

// @Summary     Scan Headers HTTP
// @Description Vérifie les headers de sécurité (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)
// @Tags        scanner
// @Produce     json
// @Param       domain query string true "Domaine à scanner"
// @Success     200 {object} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/header [get]
func handleHeader() http.HandlerFunc {
	return makeScanHandler("header", scanner.HeaderScanner{})
}

// @Summary     Scan fichiers sensibles
// @Description Détecte les fichiers sensibles exposés (.env, .git/config, wp-config.php, etc.)
// @Tags        scanner
// @Produce     json
// @Param       domain query string true "Domaine à scanner"
// @Success     200 {object} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/sensitive [get]
func handleSensitive() http.HandlerFunc {
	return makeScanHandler("sensitive", scanner.SensitiveScanner{})
}

// @Summary     Scan sous-domaines
// @Description Énumère les sous-domaines via Certificate Transparency (crt.sh)
// @Tags        scanner
// @Produce     json
// @Param       domain query string true "Domaine à scanner"
// @Success     200 {object} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/subdomain [get]
func handleSubdomain() http.HandlerFunc {
	return makeScanHandler("subdomain", scanner.SubdomainScanner{})
}

// @Summary     All Scan
// @Description Lance les 5 scanners en parallèle via goroutines
// @Tags        scanner
// @Produce     json
// @Param 		domain query string true "Domaine à scanner"
// @Success     200 {array} ScanResult
// @Failure     400 {string} string "paramètre 'domain' requis"
// @Failure     500 {string} string "erreur serveur"
// @Router      /scan/all [get]
func (s *Server) handleAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")

		if domain == "" {
			http.Error(w, "paramètre 'domain' requis", http.StatusBadRequest)
			return
		}

		// Channel pour recevoir les résultats des goroutines
		// Chaque goroutine y envoie un ScanResult quand elle a fini
		ch := make(chan ScanResult)

		// Lance une goroutine par scanner — exécution en parallèle
		// sc est passé en paramètre pour éviter les problèmes de closure
		for _, sc := range s.Scanners {
			go func(sc scanner.Scanner) {
				result, err := sc.Scan(domain)
				// Dans une goroutine, on ne peut pas faire http.Error (pas accès à w)
				// On met le message d'erreur dans Result à la place
				if err != nil {
					log.Println(err)
					result = "erreur interne du serveur"
				}

				ch <- ScanResult{
					Scanner: sc.Name(),
					Domain:  domain,
					Result:  result,
				}
			}(sc)
		}

		// Collecte les résultats — <-ch bloque jusqu'à recevoir un résultat
		// On itère autant de fois qu'il y a de scanners
		var results []ScanResult

		for i := 0; i < len(s.Scanners); i++ {
			result := <-ch
			results = append(results, result)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Println(err)
			http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
			return
		}

	}
}

// corsMiddleware — middleware qui ajoute les headers CORS à chaque réponse
// Intercepte les requêtes avant qu'elles n'atteignent le handler final
// Les requêtes OPTIONS (preflight) sont traitées directement avec un 200
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Requête preflight (OPTIONS) — le navigateur demande la permission avant le vrai appel
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Start enregistre les routes HTTP et démarre le serveur
// Toutes les routes sont enregistrées AVANT ListenAndServe (qui est bloquant)
func (s *Server) Start() {

	// Swagger UI — documentation interactive de l'API sur /swagger/index.html
	http.HandleFunc("/swagger/", httpSwagger.Handler())

	http.HandleFunc("/health", handleHealth())

	http.HandleFunc("/scan/dns", handleDNS())

	http.HandleFunc("/scan/ssl", handleSSL())

	http.HandleFunc("/scan/header", handleHeader())

	http.HandleFunc("/scan/sensitive", handleSensitive())

	http.HandleFunc("/scan/subdomain", handleSubdomain())

	http.HandleFunc("/scan/all", s.handleAll())

	// Démarrage du serveur — ListenAndServe est bloquant
	// Le programme reste ici et écoute les connexions entrantes
	fmt.Println("Serveur démarré sur le port :", s.Port)

	// fmt.Sprintf(":%d", s.Port) convertit l'int en string formatée (ex: 8082 → ":8082")
	// nil = utilise le routeur par défaut (celui où on a enregistré les routes avec HandleFunc)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.Port),
		corsMiddleware(http.DefaultServeMux),
	)

	// Si ListenAndServe retourne, c'est qu'il y a eu une erreur (ex: port déjà pris)
	// log.Fatal affiche l'erreur ET arrête le programme (exit code 1)
	if err != nil {
		log.Fatal(err)
	}
}

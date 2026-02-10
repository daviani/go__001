package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/daviani/go__001/internal/scanner"
)

type Server struct {
	Port     int
	Scanners []scanner.Scanner
}

func (s *Server) Start() {

	http.HandleFunc(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			w.Write([]byte("Status OK"))
		},
	)

	http.HandleFunc(
		"/scan/dns",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			domain := r.URL.Query().Get("domain")

			result, err := scanner.DNSScanner{}.Scan(domain)

			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		},
	)

	http.HandleFunc(
		"/scan/header",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			domain := r.URL.Query().Get("domain")

			result, err := scanner.HeaderScanner{}.Scan(domain)

			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		},
	)

	http.HandleFunc(
		"/scan/sensitive",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			domain := r.URL.Query().Get("domain")

			result, err := scanner.SensitiveScanner{}.Scan(domain)

			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		},
	)

	http.HandleFunc(
		"/scan/ssl",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			domain := r.URL.Query().Get("domain")

			result, err := scanner.SSLScanner{}.Scan(domain)

			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		},
	)

	http.HandleFunc(
		"/scan/subdomain",
		func(
			w http.ResponseWriter,
			r *http.Request) {
			domain := r.URL.Query().Get("domain")

			result, err := scanner.SubdomainScanner{}.Scan(domain)

			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			w.Write([]byte(result))
		},
	)

	fmt.Println("Serveur démarré sur le port :", s.Port)

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.Port),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}
}

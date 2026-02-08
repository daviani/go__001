package scanner

import (
	"fmt"
	"net/http"
)

// HeaderScanner - Scanner pour les headers HTTP de sécurité
type HeaderScanner struct{}

// Name retourne l'identifiant du scanner Headers
func (h HeaderScanner) Name() string { return "header" }

// Scan effectue une requête HTTP et récupère les headers de sécurité
// Vérifie HSTS, CSP et X-Frame-Options (protection contre le clickjacking)
func (h HeaderScanner) Scan(domain string) (string, error) {
	// http.Get effectue une requête GET - on ajoute https:// car domain = "daviani.dev"
	resp, err := http.Get("https://" + domain)

	if err != nil {
		return "", fmt.Errorf("erreur de header: %w", err)
	}

	// defer ferme le body à la fin de la fonction (libère les ressources)
	defer func() { _ = resp.Body.Close() }()

	// resp.Header est une map[string][]string contenant tous les headers HTTP
	headers := resp.Header

	// headers.Get("key") retourne la valeur du header ou "" si absent
	// HSTS : force HTTPS | CSP : politique de sécurité | X-Frame-Options : anti-clickjacking
	return fmt.Sprintf("HSTS: %s | CSP: %s | X-Frame-Options: %s",
		headers.Get("Strict-Transport-Security"),
		headers.Get("Content-Security-Policy"),
		headers.Get("X-Frame-Options"),
	), nil
}

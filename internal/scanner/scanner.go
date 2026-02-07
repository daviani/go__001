package scanner

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

// Scanner définit le contrat que tous les scanners doivent implémenter
// Toute struct ayant ces méthodes implémente automatiquement l'interface
type Scanner interface {
	Scan(domain string) string
	Name() string
}

// DNSScanner - Scanner pour la résolution DNS (records A et AAAA)
type DNSScanner struct{}

// SSLScanner - Scanner pour les certificats SSL/TLS
type SSLScanner struct{}

// HeaderScanner - Scanner pour les headers HTTP de sécurité
type HeaderScanner struct{}

// Name retourne l'identifiant du scanner DNS
func (d DNSScanner) Name() string { return "dns" }

// Name retourne l'identifiant du scanner SSL
func (s SSLScanner) Name() string { return "ssl" }

// Name retourne l'identifiant du scanner Headers
func (h HeaderScanner) Name() string { return "header" }

// Scan effectue une résolution DNS et retourne les adresses IP du domaine
// Utilise net.LookupIP qui résout les records A (IPv4) et AAAA (IPv6)
func (d DNSScanner) Scan(domain string) string {
	// LookupIP retourne une slice d'IPs et une erreur potentielle
	ips, err := net.LookupIP(domain)

	// Gestion d'erreur Go : on vérifie si err n'est pas nil
	if err != nil {
		return "Erreur de DNS" + err.Error()
	}

	// Construction du résultat en parcourant la slice d'IPs
	result := "Ips pour " + domain + ": "
	for _, ip := range ips {
		result += ip.String() + " "
	}
	return result
}

// Scan établit une connexion TLS et récupère les infos du certificat
// Utilise crypto/tls pour une connexion sécurisée native (pas de curl/openssl)
func (s SSLScanner) Scan(domain string) string {
	// tls.Dial ouvre une connexion TLS sur le port 443
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return "Erreur SSL: " + err.Error()
	}

	// defer garantit que la connexion sera fermée à la fin de la fonction
	// _ = ignore l'erreur de Close() volontairement
	defer func() { _ = conn.Close() }()

	// Récupère le premier certificat de la chaîne (celui du domaine)
	cert := conn.ConnectionState().PeerCertificates[0]

	// Sprintf formate les infos du certificat en une string lisible
	// Format date : "02/01/2006" = jour/mois/année (format Go spécifique)
	return fmt.Sprintf("Domaine: %s | Émetteur: %s | Expire: %s",
		cert.Subject.CommonName,
		cert.Issuer.Organization[0],
		cert.NotAfter.Format("02/01/2006"))
}

// Scan effectue une requête HTTP et récupère les headers de sécurité
// Vérifie HSTS, CSP et X-Frame-Options (protection contre le clickjacking)
func (h HeaderScanner) Scan(domain string) string {
	// http.Get effectue une requête GET - on ajoute https:// car domain = "daviani.dev"
	resp, err := http.Get("https://" + domain)

	if err != nil {
		return "Erreur: " + err.Error()
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
	)
}

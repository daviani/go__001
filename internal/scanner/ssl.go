package scanner

import (
	"crypto/tls"
	"fmt"
)

// SSLScanner - Scanner pour les certificats SSL/TLS
type SSLScanner struct{}

// Name retourne l'identifiant du scanner SSL
func (s SSLScanner) Name() string { return "ssl" }

// Scan établit une connexion TLS et récupère les infos du certificat
// Utilise crypto/tls pour une connexion sécurisée native (pas de curl/openssl)
func (s SSLScanner) Scan(domain string) (string, error) {
	// tls.Dial ouvre une connexion TLS sur le port 443
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return "", fmt.Errorf("erreur SSL: %w", err)
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
		cert.NotAfter.Format("02/01/2006")), nil
}

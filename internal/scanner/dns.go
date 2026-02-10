package scanner

import (
	"fmt"
	"net"
)

// DNSScanner - Scanner pour la résolution DNS (records A et AAAA)
type DNSScanner struct{}

// Name retourne l'identifiant du scanner DNS
func (d DNSScanner) Name() string { return "dns" }

// Scan effectue une résolution DNS complète du domaine
// Résout les records A/AAAA (IPs), MX (serveurs mail), NS (nameservers) et TXT (SPF, DMARC...)
func (d DNSScanner) Scan(domain string) (string, error) {

	// --- Records A et AAAA (adresses IP) ---
	// LookupIP retourne une slice de net.IP (IPv4 + IPv6)
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", fmt.Errorf("erreur de DNS: %w", err)
	}

	// Construction des résultats par type de record
	resultIP := "IPs pour " + domain + ": \n"
	resultMX := "MXs pour " + domain + ": \n"
	resultNS := "NSs pour " + domain + ": \n"
	resultTXT := "TXTs pour " + domain + ": \n"

	// --- Records MX (serveurs mail) ---
	// LookupMX retourne []*net.MX — chaque MX a un champ .Host (string) et .Pref (priorité)
	mxs, err := net.LookupMX(domain)
	if err != nil {
		resultMX = "MX : erreur de résolution"
	}

	// --- Records NS (nameservers) ---
	// LookupNS retourne []*net.NS — chaque NS a un champ .Host (string)
	nss, err := net.LookupNS(domain)
	if err != nil {
		resultNS = "NS : erreur de résolution"
	}

	// --- Records TXT (SPF, vérification domaine...) ---
	// LookupTXT retourne directement []string — pas besoin de .Host ou .String()
	txts, err := net.LookupTXT(domain)
	if err != nil {
		resultTXT = "TXT : erreur de résolution"
	}

	// ip.String() convertit net.IP en string lisible (ex: "188.114.96.2")
	for _, ip := range ips {
		resultIP += ip.String() + "\n"
	}

	// mx.Host est un champ string de la struct net.MX (ex: "mx01.mail.icloud.com.")
	for _, mx := range mxs {
		resultMX += mx.Host + "\n"
	}

	// ns.Host est un champ string de la struct net.NS (ex: "fish.ns.cloudflare.com.")
	for _, ns := range nss {
		resultNS += ns.Host + "\n"
	}

	// txt est déjà une string, pas besoin de conversion
	for _, txt := range txts {
		resultTXT += txt + "\n"
	}

	return resultIP + resultMX + resultNS + resultTXT, nil
}

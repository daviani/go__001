package scanner

import "net"

// DNSScanner - Scanner pour la résolution DNS (records A et AAAA)
type DNSScanner struct{}

// Name retourne l'identifiant du scanner DNS
func (d DNSScanner) Name() string { return "dns" }

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

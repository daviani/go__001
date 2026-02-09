package scanner

import "testing"

// TestSSLScanner_Name vérifie que le scanner retourne le bon identifiant
func TestSSLScanner_Name(t *testing.T) {
	result := SSLScanner{}.Name()

	if result != "ssl" {
		t.Errorf("got %s, want ssl", result)
	}
}

// TestSSLScanner_Scan — Happy path : vérifie qu'un domaine valide retourne les infos du certificat
// SSLScanner utilise tls.Dial sur le port 443 pour récupérer le certificat x509
func TestSSLScanner_Scan(t *testing.T) {
	// tls.Dial("tcp", "google.com:443", nil) en interne
	result, err := SSLScanner{}.Scan("google.com")

	// La connexion TLS doit réussir — Google a un certificat valide
	if err != nil {
		t.Fatal(err)
	}

	// Le résultat contient "Domaine: ... | Émetteur: ... | Expire: ..."
	if result == "" {
		t.Errorf("got empty result, want certificate info")
	}
}

// TestSSLScanner_Scan_InvalidDomain — Error path : un domaine invalide fait échouer tls.Dial
// "false_url" → résolution DNS échoue → tls.Dial retourne une erreur immédiatement
func TestSSLScanner_Scan_InvalidDomain(t *testing.T) {
	result, err := SSLScanner{}.Scan("false_url")

	// tls.Dial échoue car le domaine n'existe pas (pas de handshake TLS possible)
	if err == nil {
		t.Errorf("expected error for invalid domain, got nil")
	}

	// Contrat : en cas d'erreur, result doit être "" (string vide)
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}

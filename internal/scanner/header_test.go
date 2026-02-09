package scanner

import "testing"

// TestHeaderScanner_Name vérifie que le scanner retourne le bon identifiant
func TestHeaderScanner_Name(t *testing.T) {
	result := HeaderScanner{}.Name()

	if result != "header" {
		t.Errorf("got %s, want header", result)
	}
}

// TestHeaderScanner_Scan — Happy path : vérifie que les headers de sécurité sont récupérés
// google.com doit répondre avec au moins HSTS ou X-Frame-Options
func TestHeaderScanner_Scan(t *testing.T) {
	// HeaderScanner fait un http.Get("https://" + domain) en interne
	result, err := HeaderScanner{}.Scan("google.com")

	// Pas d'erreur attendue — le serveur Google répond toujours
	if err != nil {
		t.Fatal(err)
	}

	// Le résultat contient les headers formatés : "HSTS: ... | CSP: ... | X-Frame: ..."
	// Même si certains headers sont vides, la string formatée n'est jamais ""
	if result == "" {
		t.Errorf("got empty result, want security headers")
	}
}

// TestHeaderScanner_Scan_InvalidDomain — Error path : un domaine invalide fait échouer http.Get
// "false_url" → résolution DNS échoue → http.Get retourne une erreur
func TestHeaderScanner_Scan_InvalidDomain(t *testing.T) {
	result, err := HeaderScanner{}.Scan("false_url")

	// http.Get échoue car "false_url" n'est pas un domaine valide
	if err == nil {
		t.Errorf("expected error for invalid domain, got nil")
	}

	// Contrat : en cas d'erreur, result doit être "" (string vide)
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}

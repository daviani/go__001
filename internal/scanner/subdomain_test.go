package scanner

import "testing"

// TestSubdomainScanner_Name vérifie que le scanner retourne le bon identifiant
func TestSubdomainScanner_Name(t *testing.T) {
	result := SubdomainScanner{}.Name()

	if result != "subdomain" {
		t.Errorf("got %s, want subdomain", result)
	}
}

// TestSubdomainScanner_Scan — Happy path : vérifie que crt.sh retourne des sous-domaines
// Ce test est plus lent (~3-20s) car il appelle l'API externe crt.sh
func TestSubdomainScanner_Scan(t *testing.T) {
	// crt.sh interroge les Certificate Transparency logs pour trouver les sous-domaines
	result, err := SubdomainScanner{}.Scan("google.com")

	if err != nil {
		t.Fatal(err)
	}

	// Le résultat doit contenir "Sous domaines: " suivi d'au moins un sous-domaine
	if result == "" {
		t.Errorf("got empty result, want subdomains list")
	}
}

// TestSubdomainScanner_InvalidDomain — Error path : un caractère null (\x00) fait échouer http.Get
// Note : un domaine bidon (ex: "false_url") ne cause PAS d'erreur car crt.sh répond avec []
// On utilise \x00 (caractère de contrôle) car http.Get refuse les URL avec des caractères invalides
func TestSubdomainScanner_InvalidDomain(t *testing.T) {
	result, err := SubdomainScanner{}.Scan("\x00")

	// http.Get refuse les URL contenant des caractères de contrôle → erreur immédiate
	if err == nil {
		t.Errorf("expected error for invalid domain, got nil")
	}

	// Contrat : en cas d'erreur, result doit être "" (string vide)
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}

package scanner

import "testing"

// TestDNSScanner_Name vérifie que le scanner retourne le bon identifiant
// Convention Go : TestNomStruct_Methode — permet de cibler un test avec -run
func TestDNSScanner_Name(t *testing.T) {
	// Crée une instance de DNSScanner et appelle Name()
	// En Go, on peut chaîner struct{} + méthode sur une seule ligne
	result := DNSScanner{}.Name()

	// Errorf = le test échoue mais continue (vérifie les autres assertions)
	// %s = placeholder pour une string (interpolation Go)
	if result != "dns" {
		t.Errorf("got %s, want dns", result)
	}
}

// TestDNSScanner_Scan — Happy path : un domaine valide doit retourner des records DNS
// Note : ce test fait un VRAI appel réseau → dépend de la connexion internet
func TestDNSScanner_Scan(t *testing.T) {
	// Scan retourne (string, error) — on capte les deux valeurs
	result, err := DNSScanner{}.Scan("google.com")

	// Fatal = arrête le test immédiatement, pas la peine de vérifier result
	// On vérifie err EN PREMIER : si le scan a échoué, result est "" (inutile à tester)
	if err != nil {
		t.Fatal(err)
	}

	// Si pas d'erreur, on vérifie que le résultat contient quelque chose
	// Errorf (pas Fatal) : non bloquant, le test continue après
	if result == "" {
		t.Errorf("got empty result, want DNS records")
	}
}

// TestDNSScanner_Scan_InvalidDomain — Error path : un domaine invalide doit retourner une erreur
// "false_url" n'existe pas → net.LookupIP échoue → err != nil attendu
func TestDNSScanner_Scan_InvalidDomain(t *testing.T) {
	result, err := DNSScanner{}.Scan("false_url")

	// Ici on VEUT une erreur — si err est nil, le test a un problème
	// Errorf (pas Fatal) : on veut aussi vérifier result ensuite
	if err == nil {
		t.Errorf("expected error for invalid domain, got nil")
	}

	// Quand le scan échoue, result doit être "" (string vide)
	// C'est le contrat de notre interface : erreur → pas de résultat
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}

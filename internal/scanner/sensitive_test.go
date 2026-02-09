package scanner

import "testing"

// TestSensitiveScanner_Name vérifie que le scanner retourne le bon identifiant
func TestSensitiveScanner_Name(t *testing.T) {
	result := SensitiveScanner{}.Name()

	if result != "sensitive" {
		t.Errorf("got %s, want sensitive", result)
	}
}

// TestSensitiveScanner_Scan — Happy path : vérifie que le scan de fichiers sensibles fonctionne
// google.com ne devrait pas exposer de fichiers sensibles → résultat = "Aucun fichier sensible trouvé"
func TestSensitiveScanner_Scan(t *testing.T) {
	// SensitiveScanner teste 6 chemins (.git/config, .env, .htaccess, robots.txt, sitemap.xml, wp-config.php)
	// via http.Get sur chaque chemin, vérifie si status == 200
	result, err := SensitiveScanner{}.Scan("google.com")

	// Les requêtes HTTP doivent aboutir (même si le fichier n'existe pas → 404)
	if err != nil {
		t.Fatal(err)
	}

	// google.com retourne "Aucun fichier sensible trouvé" ou la liste des fichiers exposés
	// Dans les deux cas, result n'est jamais vide
	if result == "" {
		t.Errorf("got empty result, want scan report")
	}
}

// TestSensitiveScanner_Scan_InvalidDomain — Error path : un domaine invalide fait échouer http.Get
// "false_url" → résolution DNS échoue dès le premier http.Get → err != nil
func TestSensitiveScanner_Scan_InvalidDomain(t *testing.T) {
	result, err := SensitiveScanner{}.Scan("false_url")

	// http.Get échoue car "false_url" ne résout pas en IP
	if err == nil {
		t.Errorf("expected error for invalid domain, got nil")
	}

	// Contrat : en cas d'erreur, result doit être "" (string vide)
	if result != "" {
		t.Errorf("expected empty result, got %s", result)
	}
}

package scanner

import (
	"fmt"
	"net/http"
)

// SensitiveScanner - Scanner pour la détection de fichiers sensibles exposés publiquement
type SensitiveScanner struct{}

// Name retourne l'identifiant du scanner Sensitive
func (d SensitiveScanner) Name() string { return "sensitive" }

// Scan teste une liste de chemins sensibles via HTTP GET
// Un status 200 signifie que le fichier est accessible publiquement → alerte de sécurité
func (d SensitiveScanner) Scan(domain string) (string, error) {

	// Liste des chemins sensibles à tester
	// .git/config → exposition du repo Git
	// .env → secrets (clés API, mots de passe)
	// .htaccess → configuration Apache
	// robots.txt / sitemap.xml → fichiers normaux mais informatifs
	// wp-config.php → configuration WordPress (accès BDD)
	paths := []string{".git/config", ".env", ".htaccess", "robots.txt", "sitemap.xml", "wp-config.php"}

	var result = ""

	for _, path := range paths {
		resp, err := http.Get("https://" + domain + "/" + path)
		if err != nil {
			return "", fmt.Errorf("erreur https: %w", err)
		}

		// StatusCode == 200 → fichier accessible publiquement
		// resp.Status contient le code + texte (ex: "200 OK", "404 Not Found")
		if resp.StatusCode == 200 {
			result += path + " → " + resp.Status + "\n"
		}

	}

	// Si aucun fichier sensible trouvé, retourner un message explicite
	if result == "" {
		return "Aucun fichier sensible trouvé", nil
	}

	return result, nil
}

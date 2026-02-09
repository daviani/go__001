package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CrtShEntry représente une entrée de la réponse JSON de l'API crt.sh
// Le tag `json:"name_value"` indique à json.Unmarshal quel champ JSON mapper
type CrtShEntry struct {
	NameValue string `json:"name_value"`
}

// SubdomainScanner - Scanner pour l'énumération de sous-domaines via Certificate Transparency
type SubdomainScanner struct{}

// Name retourne l'identifiant du scanner Subdomain
func (sb SubdomainScanner) Name() string { return "subdomain" }

// Scan interroge l'API crt.sh pour trouver tous les sous-domaines
// ayant un certificat SSL émis pour le domaine cible
func (sb SubdomainScanner) Scan(domain string) (string, error) {

	// Construction de l'URL crt.sh — %%25 = %25 encodé (wildcard %)
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erreur de subdomain: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Lit le body HTTP en entier et le désérialise en slice de CrtShEntry
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur de lecture: %w", err)
	}

	// &results : pointeur nécessaire pour que Unmarshal puisse remplir la slice
	var results []CrtShEntry
	if err = json.Unmarshal(body, &results); err != nil {
		return "", fmt.Errorf("erreur de désérialisation: %w", err)
	}

	// Dédoublonnage — map[string]bool utilisée comme Set
	unique := make(map[string]bool)
	for _, entry := range results {
		unique[entry.NameValue] = true
	}

	// Extrait les clés uniques et les joint en string
	// Équivalent JS : Object.keys(unique).join(", ")
	var keys []string
	for key := range unique {
		keys = append(keys, key)
	}

	return "Sous domaines: " + strings.Join(keys, " , "), nil
}

package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (sb SubdomainScanner) Scan(domain string) string {

	// Construction de l'URL crt.sh - %%25 = %25 encodé (wildcard %)
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	// Requête HTTP GET vers l'API crt.sh
	resp, err := http.Get(url)
	if err != nil {
		return "Erreur: " + err.Error()
	}

	// Ferme le body à la fin de la fonction (libère les ressources réseau)
	defer func() { _ = resp.Body.Close() }()

	// io.ReadAll lit tout le stream HTTP et retourne un []byte (slice d'octets)
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "Erreur: " + err.Error()
	}

	// Désérialise le JSON brut en slice de CrtShEntry
	// &results = pointeur pour que Unmarshal puisse modifier la variable directement
	var results []CrtShEntry
	err = json.Unmarshal(body, &results)

	if err != nil {
		return "Erreur JSON: " + err.Error()
	}

	// map[string]bool utilisée comme Set (dédoublonnage)
	// Les maps sont des types référence : pas besoin de & pour les modifier
	unique := make(map[string]bool)

	// Boucle 1 : parcourt les résultats JSON et ajoute chaque sous-domaine dans la map
	// Si le sous-domaine existe déjà, il est simplement écrasé (même valeur true)
	for _, entry := range results {
		unique[entry.NameValue] = true
	}

	// Boucle 2 : parcourt les clés de la map (sous-domaines uniques)
	// et construit la string de retour
	result := "Sous domaines: "

	for key := range unique {
		result += key + " , "
	}

	return result
}

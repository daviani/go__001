package main

import (
	"flag"
	"fmt"

	"github.com/daviani/go__001/internal/config"
	"github.com/daviani/go__001/internal/scanner"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}
func main() {
	// Initialisation des scanners (structs vides qui implémentent l'interface Scanner)
	dns := scanner.DNSScanner{}
	ssl := scanner.SSLScanner{}
	header := scanner.HeaderScanner{}
	subdomain := scanner.SubdomainScanner{}
	sensitive := scanner.SensitiveScanner{}

	defaultDomain := config.RequireEnv("DEFAULT_DOMAIN")
	// Slice contenant tous les scanners - on peut en ajouter autant qu'on veut
	scanners := []scanner.Scanner{dns, ssl, header, subdomain, sensitive}

	var domain string

	flag.StringVar(&domain, "domain", defaultDomain, "Domaine à scanner")
	flag.StringVar(&domain, "d", defaultDomain, "Domaine à scanner (raccourci)")

	flag.Parse()
	// Channel pour la communication entre goroutines
	// Les goroutines enverront leurs résultats ici
	ch := make(chan scanner.Result)

	// Boucle 1 : Lance une goroutine par scanner (exécution en parallèle)
	// Chaque goroutine exécute Scan() et envoie le résultat dans le channel
	for _, value := range scanners {
		go func(s scanner.Scanner) {
			result, err := s.Scan(domain)
			if err != nil {
				result = "Erreur: " + err.Error()
			}
			ch <- scanner.Result{
				Name:   s.Name(),
				Result: result,
			}
		}(value) // On passe "value" en paramètre pour éviter les problèmes de closure
	}

	// Boucle 2 : Récupère les résultats du channel
	// On itère autant de fois qu'il y a de scanners (= autant de résultats attendus)
	// <-ch bloque jusqu'à ce qu'une goroutine envoie un résultat

	// Stocke les résultats dans une map[nom]résultat pour un affichage ordonné
	// Sans la map, l'ordre dépendrait de quelle goroutine finit en premier
	results := make(map[string]string)

	for i := 0; i < len(scanners); i++ {
		sr := <-ch
		results[sr.Name] = sr.Result
	}
	// Définition des sections du rapport avec struct anonyme
	// Permet d'afficher les résultats dans un ordre défini (pas l'ordre aléatoire des goroutines)
	sections := []struct {
		Title string
		Key   string
	}{
		{"DNS", "dns"},
		{"SSL", "ssl"},
		{"HEADERS", "header"},
		{"SOUS-DOMAINES", "subdomain"},
		{"FICHIERS SENSIBLES", "sensitive"},
	}

	// Affichage du rapport — header unique avant la boucle
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("── Rapport de sécurité — " + domain + " ──────────────")
	fmt.Println("══════════════════════════════════════════")

	for _, section := range sections {
		printSection(section.Title, results[section.Key])
	}

	fmt.Println("═════════════ END ═════════════════")
}

// printSection affiche une section du rapport avec un titre formaté
func printSection(title string, content string) {
	fmt.Println("── " + title + " ──────────────────────────────")
	fmt.Println(content)
}

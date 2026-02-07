package main

import (
	"fmt"

	"github.com/daviani/go__001/internal/scanner"
)

func main() {
	// Initialisation des scanners (structs vides qui implémentent l'interface Scanner)
	dns := scanner.DNSScanner{}
	ssl := scanner.SSLScanner{}
	header := scanner.HeaderScanner{}

	// Slice contenant tous les scanners - on peut en ajouter autant qu'on veut
	scanners := []scanner.Scanner{dns, ssl, header}

	// Channel pour la communication entre goroutines
	// Les goroutines enverront leurs résultats ici
	ch := make(chan string)

	// Boucle 1 : Lance une goroutine par scanner (exécution en parallèle)
	// Chaque goroutine exécute Scan() et envoie le résultat dans le channel
	for _, value := range scanners {
		go func(s scanner.Scanner) {
			result := s.Scan("daviani.dev")
			ch <- result // Envoie le résultat dans le channel
		}(value) // On passe "value" en paramètre pour éviter les problèmes de closure
	}

	// Boucle 2 : Récupère les résultats du channel
	// On itère autant de fois qu'il y a de scanners (= autant de résultats attendus)
	// <-ch bloque jusqu'à ce qu'une goroutine envoie un résultat
	for i := 0; i < len(scanners); i++ {
		result := <-ch // Reçoit un résultat du channel
		fmt.Println(result)
	}
}

package scanner

// Scanner définit le contrat que tous les scanners doivent implémenter
// Toute struct ayant ces méthodes implémente automatiquement l'interface
type Scanner interface {
	Scan(domain string) (string, error)
	Name() string
}

// Result transporte le résultat d'un scan via le channel
// Permet d'associer le nom du scanner à son résultat (goroutines → main)
type Result struct {
	Name   string
	Result string
}

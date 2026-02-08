package scanner

// Scanner définit le contrat que tous les scanners doivent implémenter
// Toute struct ayant ces méthodes implémente automatiquement l'interface
type Scanner interface {
	Scan(domain string) string
	Name() string
}

// ScannerResult transporte le résultat d'un scan via le channel
// Permet d'associer le nom du scanner à son résultat (goroutines → main)
type ScannerResult struct {
	Name   string
	Result string
}

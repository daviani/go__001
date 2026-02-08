package scanner

// Scanner définit le contrat que tous les scanners doivent implémenter
// Toute struct ayant ces méthodes implémente automatiquement l'interface
type Scanner interface {
	Scan(domain string) string
	Name() string
}

package config

import (
	"fmt"
	"os"
)

// GetEnv récupère une variable d'environnement optionnelle.
// Si la variable n'est pas définie ou vide, retourne la valeur fallback.
// Utiliser pour : configuration applicative (port, timeout, domaine par défaut...)
func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// RequireEnv récupère une variable d'environnement obligatoire.
// Si la variable n'est pas définie ou vide, le programme crash immédiatement (panic).
// Utiliser pour : secrets et configurations critiques (DATABASE_URL, API_KEY...)
func RequireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Variable d'environnement requise : %s", key))
	}

	return value
}

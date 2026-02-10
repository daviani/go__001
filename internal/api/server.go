package api

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/daviani/go__001/docs"          // Blank import — enregistre la spec Swagger au démarrage via init()
	httpSwagger "github.com/swaggo/http-swagger" // Middleware servant l'interface Swagger UI
)

// Start enregistre les routes HTTP et démarre le serveur
// Toutes les routes sont enregistrées AVANT ListenAndServe (qui est bloquant)
func (s *Server) Start() {

	// Swagger UI — documentation interactive de l'API sur /swagger/index.html
	http.HandleFunc("/swagger/", httpSwagger.Handler())

	http.HandleFunc("/health", handleHealth())

	http.HandleFunc("/scan/dns", handleDNS())

	http.HandleFunc("/scan/ssl", handleSSL())

	http.HandleFunc("/scan/header", handleHeader())

	http.HandleFunc("/scan/sensitive", handleSensitive())

	http.HandleFunc("/scan/subdomain", handleSubdomain())

	http.HandleFunc("/scan/all", s.handleAll())

	// Démarrage du serveur — ListenAndServe est bloquant
	// Le programme reste ici et écoute les connexions entrantes
	fmt.Println("Serveur démarré sur le port :", s.Port)

	// fmt.Sprintf(":%d", s.Port) convertit l'int en string formatée (ex: 8082 → ":8082")
	// nil = utilise le routeur par défaut (celui où on a enregistré les routes avec HandleFunc)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.Port),
		corsMiddleware(http.DefaultServeMux),
	)

	// Si ListenAndServe retourne, c'est qu'il y a eu une erreur (ex: port déjà pris)
	// log.Fatal affiche l'erreur ET arrête le programme (exit code 1)
	if err != nil {
		log.Fatal(err)
	}
}

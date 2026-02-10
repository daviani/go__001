package api

import "net/http"

// corsMiddleware — middleware qui ajoute les headers CORS à chaque réponse
// Intercepte les requêtes avant qu'elles n'atteignent le handler final
// Les requêtes OPTIONS (preflight) sont traitées directement avec un 200
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Requête preflight (OPTIONS) — le navigateur demande la permission avant le vrai appel
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

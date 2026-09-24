package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middleWareMedtricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

// func (cfg *apiConfig) middleWareResponse(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		count := fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())
// 		w.Write([]byte(count))
// 	})
// }

func (cfg *apiConfig) middleWareMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	count := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
	w.Write([]byte(count))
}

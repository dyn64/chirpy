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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	count := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(count))
}

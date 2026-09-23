package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	const rootpath = "."

	mux := http.NewServeMux()

	fserv := http.FileServer(http.Dir(rootpath + "/pub/"))

	mux.Handle("/", fserv)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting server on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}

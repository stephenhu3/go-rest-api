package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	// Https cert and key generation and usage -> removed for demo
	// generateCertKeyPEM()
	// log.Fatal( http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router))
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

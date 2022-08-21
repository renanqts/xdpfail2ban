package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/drop", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("dummy requested %s %s %s", r.RequestURI, r.Method, r.RemoteAddr)
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			log.Printf("dummy response code %d", http.StatusCreated)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
			log.Printf("dummy response code %d", http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
			log.Printf("dummy response code %d", http.StatusOK)
		}
	})

	log.Println("Starting dummy xdpdropper on port :8081")
	_ = http.ListenAndServe(":8081", nil)
}

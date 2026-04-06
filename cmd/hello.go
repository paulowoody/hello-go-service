package main

import (
	"fmt"
	"net/http"
)

var isHealthy = true

func main() {
	// Main application logic
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "Hello, World!")
		if err != nil {
			return
		}
	})

	// Health check endpoint (Liveness)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !isHealthy {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "unhealthy")
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "OK")
		if err != nil {
			return
		}
	})

	// Readiness endpoint (Readiness)
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// In a real app, you might check a database connection here
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "Ready")
		if err != nil {
			return
		}
	})

	// The Poison Pill endpoint
	http.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		isHealthy = false
		_, err := fmt.Fprint(w, "Pod is now failing health checks.")
		if err != nil {
			return
		}
	})

	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", handleHealthz)
	http.HandleFunc("/webhook", handleWebhook)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

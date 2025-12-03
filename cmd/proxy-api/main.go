package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("/etc/proxyx/web-admin"))
	http.Handle("/", fs)

	http.HandleFunc("/api/hello", func (w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("Hello from the server"))
	})

	if err := http.ListenAndServe(":5053", nil ); err != nil {
		log.Fatalf("Error starting Web Server: %v", err)
	}
}
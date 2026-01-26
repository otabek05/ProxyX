package main

import (
	"encoding/json"
	"log"
	"net/http"
)


type ApiResponse struct {
	Message string `json:"message"`
	Data any `json:"data"`
	StatusCode string `json:"status_code"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	log.Printf(
		"remote=%s | method=%s | url=%s\n",
		r.RemoteAddr,
		r.Method,
		r.URL.String(),
	)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := ApiResponse{
		Message:    "success",
		Data:       "response from 8080",
		StatusCode: "1002",
	}

	_ = json.NewEncoder(w).Encode(response)
})


	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
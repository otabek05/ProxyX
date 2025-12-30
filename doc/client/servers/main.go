package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


type Response struct {
	Message string 
	StatusCode int 
	Data any 
}

func main() {
	port := "8081"

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request arrived on: ", r.Method, r.URL.Path, r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")

		message := fmt.Sprintf("server is running in port:%s", port)
		response := &Response{
			Message: message,
			StatusCode: 200,
			Data: "something from the web",
		}


		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})


	server :=  &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Println("Server started on ", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
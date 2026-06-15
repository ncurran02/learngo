package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Unable to load .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic... pretend something went wrong!")
	})

	handler := RecoverMiddleware{
		Next: mux,
	}

	log.Printf("Listening on :8080")
	http.ListenAndServe(":8080", handler)
}

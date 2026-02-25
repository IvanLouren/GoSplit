package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IvanLouren/GoSplit/pkg/database"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	database.Connect()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

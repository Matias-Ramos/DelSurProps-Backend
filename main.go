package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/crud"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	//******************************************
	// SV, Credentials & DB Initialization.
	sv := chi.NewRouter()
	dbAuth, mapError := godotenv.Read(".env")
	if mapError != nil {
		fmt.Printf("Error loading .env into map[string]string - %s", mapError)
		return
	}
	db, sqlErr := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", dbAuth["PGS_USER"], dbAuth["PGS_PWD"], dbAuth["PGS_DB_NAME"]))
	if sqlErr != nil {
		fmt.Printf("DB initialization failed - %s", sqlErr)
		return
	}
	defer db.Close()

	//******************************************
	// CORS Middleware to open API-React traffic.
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your React app's URL
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	sv.Use(cors.Handler)

	//******************************************
	// Handlers.
	categoryHandler := func(w http.ResponseWriter, r *http.Request) {
		crud.GetDBdata(w, r, db)
	}
	sv.Get("/api/{category}", categoryHandler)
	sv.Post("/admin/post/{category}", crud.PostData)

	//******************************************
	// Turning on the server.
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/auth"
	"github.com/Matias-Ramos/Inmobiliaria-backend-go/crud"
	"github.com/Matias-Ramos/Inmobiliaria-backend-go/logs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var apiKey, jwtSecret, pgsUser, pgsPwd, pgsDbName string

func init() {
	rand.Seed(time.Now().UnixNano())

	// ****************************
	// dotenv config.
	dotEnv, mapError := godotenv.Read(".env")
	if mapError != nil {
		fmt.Printf("Error loading .env into map[string]string - %s", mapError)
		return
	}
	apiKey, jwtSecret = dotEnv["API_KEY"], dotEnv["JWT_SECRET"]
	pgsUser, pgsPwd, pgsDbName = dotEnv["PGS_USER"], dotEnv["PGS_PWD"], dotEnv["PGS_DB_NAME"]

}

func main() {

	//******************************************
	// SV Init

	sv := chi.NewRouter()
	sv.Use(middleware.Recoverer)

	//******************************************
	// Logging Init.

	logFile := logs.OpenLogFile()
	defer logFile.Close()
	log.SetOutput(logFile)

	//******************************************
	// DB Init

	db, sqlErr := sql.Open("postgres", 
	fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", pgsUser, pgsPwd, pgsDbName))
	if sqlErr != nil {
		log.Printf("DB initialization failed - %s", sqlErr)
		return
	}
	defer db.Close()

	//******************************************
	// CORS Middleware to open API-React traffic.

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	sv.Use(cors.Handler)

	//******************************************
	// Handlers.

	sv.Get("/api/{category}", crud.GetDBdata(db))

	sv.Post("/admin/jwt", auth.GetJwt(apiKey, jwtSecret))
	sv.With(auth.ValidateJwt(jwtSecret)).Post("/admin/post/{category}", http.HandlerFunc(crud.PostData(db)))

	// sv.Post("/admin/post/{category}", crud.PostData)

	//******************************************
	// Turning on the server.

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

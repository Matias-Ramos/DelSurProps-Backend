package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func fillBuildingDetails(category string, rows *sql.Rows) (interface{}, error) {
	switch category {
	case "Alquileres":
		buildingObj := &RentBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Location,
			&buildingObj.Price,
			&buildingObj.Env,
			&buildingObj.Bedrooms,
			&buildingObj.Bathrooms,
			&buildingObj.Garages,
			pq.Array(&buildingObj.Images))
		return buildingObj, err
	case "Ventas":
		buildingObj := &SalesBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Location,
			&buildingObj.Price,
			&buildingObj.Env,
			&buildingObj.Bedrooms,
			&buildingObj.Bathrooms,
			&buildingObj.Garages,
			&buildingObj.Covered_surface,
			&buildingObj.Total_surface,
			pq.Array(&buildingObj.Images))
		return buildingObj, err
	case "Emprendimientos":
		buildingObj := &VentureBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Location,
			&buildingObj.Price,
			&buildingObj.Env,
			&buildingObj.Bedrooms,
			&buildingObj.Bathrooms,
			&buildingObj.Garages,
			&buildingObj.Covered_surface,
			&buildingObj.Total_surface,
			&buildingObj.Pozo,
			&buildingObj.In_progress,
			pq.Array(&buildingObj.Images))
		return buildingObj, err
	default:
		return nil, fmt.Errorf("unsupported category: %s", category)
	}
}
func generateSQLquery(category string, urlQyParams url.Values) (string, []interface{}) {
	// Building the SQL query
	// (this way to query prevents SQL injection vulnerabilities)
	query := fmt.Sprintf(`SELECT * FROM public."%s"`, category)
	args := []interface{}{}
	conditions := []string{}

	for fieldKey, fieldValue := range urlQyParams {
		switch fieldKey {
		case "location":
			conditions = append(conditions, "location ILIKE $"+strconv.Itoa(len(args)+1))
			args = append(args, "'%"+fieldValue[0]+"%'")

		case "price_init":
			conditions = append(conditions, "price >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "price_limit":
			conditions = append(conditions, "price <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "env_init":
			conditions = append(conditions, "env >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "env_limit":
			conditions = append(conditions, "env <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "bedroom_init":
			conditions = append(conditions, "bedrooms >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "bedroom_limit":
			conditions = append(conditions, "bedrooms <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "bathroom_init":
			conditions = append(conditions, "bathrooms >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "bathroom_limit":
			conditions = append(conditions, "bathrooms <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "garage_init":
			conditions = append(conditions, "garages >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "garage_limit":
			conditions = append(conditions, "garages <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "total_surface_init":
			conditions = append(conditions, "total_surface >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "total_surface_limit":
			conditions = append(conditions, "total_surface <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "covered_surface_init":
			conditions = append(conditions, "covered_surface >= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])
		case "covered_surface_limit":
			conditions = append(conditions, "covered_surface <= $"+strconv.Itoa(len(args)+1))
			args = append(args, fieldValue[0])

		case "building_status":
			switch fieldValue[0] {
			case "in_progress":
				conditions = append(conditions, "in_progress = $"+strconv.Itoa(len(args)+1))
				args = append(args, true)
			case "pozo":
				conditions = append(conditions, "pozo = $"+strconv.Itoa(len(args)+1))
				args = append(args, true)
			case "in_progress-or-pozo":
				conditions = append(conditions, "in_progress = $"+strconv.Itoa(len(args)+1)+" OR pozo = $"+strconv.Itoa(len(args)+2))
				args = append(args, true, true)
			}
		}
	}
	// Add the WHERE clause if there are conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query, args
}
func getDBdata(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	//***************************************
	// DB data gathering through SQL querying.
	category := r.URL.Path[1:]
	urlQyParams := r.URL.Query()
	query, args := generateSQLquery(category, urlQyParams)
	rows, err := db.Query(query, args...) // ($x from "query") are replaced by ('value' from "args")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//********************************
	// Slice of structs initialization.
	var buildings []interface{}
	for rows.Next() {
		newBuilding, err := fillBuildingDetails(category, rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buildings = append(buildings, newBuilding)
	}

	//**********************************
	// Convertion from Go slice to JSON.
	jsonData, err := json.Marshal(buildings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//**********************************
	// Sending the data to the requester.
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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
	db, sqlErr := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", dbAuth["USER"], dbAuth["PWD"], dbAuth["DB_NAME"]))
	if sqlErr != nil {
		fmt.Printf("DB initialization failed - %s", sqlErr)
		return
	}
	defer db.Close()

	//******************************************
	// CORS Middleware to open API-React traffic.
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your React app's URL
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	sv.Use(cors.Handler)

	//******************************************
	// Handlers.
	sv.Get("/Emprendimientos", func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db)
	})
	sv.Get("/Ventas", func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db)
	})
	sv.Get("/Alquileres", func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db)
	})

	//******************************************
	// Turning on the server.
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

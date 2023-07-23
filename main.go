package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/env/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/lib/pq"
)

func fillBuildingDetails(category string, rows *sql.Rows) (interface{}, error) {

	switch category {
	case "Alquileres":
		buildingObj := &RentBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Ubicacion,
			&buildingObj.Precio,
			&buildingObj.Ambientes,
			&buildingObj.Dormitorios,
			&buildingObj.Banios,
			&buildingObj.Garages,
			pq.Array(&buildingObj.Imagenes))
		return buildingObj, err
	case "Ventas":
		buildingObj := &SalesBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Ubicacion,
			&buildingObj.Precio,
			&buildingObj.Ambientes,
			&buildingObj.Dormitorios,
			&buildingObj.Banios,
			&buildingObj.Garages,
			&buildingObj.Superficie_cubierta,
			&buildingObj.Superficie_total,
			&buildingObj.Imagenes)
		return buildingObj, err
	case "Emprendimientos":
		buildingObj := &VentureBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Ubicacion,
			&buildingObj.Precio,
			&buildingObj.Ambientes,
			&buildingObj.Dormitorios,
			&buildingObj.Banios,
			&buildingObj.Garages,
			&buildingObj.Superficie_cubierta,
			&buildingObj.Superficie_total,
			&buildingObj.En_pozo,
			&buildingObj.En_construccion,
			&buildingObj.Imagenes)
		return buildingObj, err
	default:
		return nil, fmt.Errorf("unsupported category: %s", category)
	}
}

func getDBdata(w http.ResponseWriter, r *http.Request, db *sql.DB, category string) {

	//***************************************
	// DB data gathering through SQL querying.
	query := fmt.Sprintf(`SELECT * FROM public."%s" ORDER BY id ASC`, category)
	rows, err := db.Query(query)
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
	config.SetEnv()
	user, pwd, db_name := config.GetEnv()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, pwd, db_name))
	if err != nil {
		log.Fatal("DB initialization - ", err)
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
		getDBdata(w, r, db, "Emprendimientos")
	})
	sv.Get("/Ventas", func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db, "Ventas")
	})
	sv.Get("/Alquileres", func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db, "Alquileres")
	})

	//******************************************
	// Turning on the server.
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

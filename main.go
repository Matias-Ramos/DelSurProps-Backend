package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

/*
InitBuilingType transmutes a *sql.Rows into a Go interface{}.
Such result will represent the building object.
*/
func initBuildingType(category string, rows *sql.Rows) (interface{}, error) {
	switch category {
	case "alquiler-inmuebles":
		buildingObj := &RentBuilding{Building: &Building{}}
		err := rows.Scan(
			&buildingObj.Id,
			&buildingObj.Location,
			&buildingObj.Price,
			&buildingObj.Env,
			&buildingObj.Bedrooms,
			&buildingObj.Bathrooms,
			&buildingObj.Garages,
			pq.Array(&buildingObj.Images),
			&buildingObj.LinkML,
			&buildingObj.LinkZonaprop,
			&buildingObj.LinkArgenprop)
		return buildingObj, err
	case "venta-inmuebles":
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
			pq.Array(&buildingObj.Images),
			&buildingObj.LinkML,
			&buildingObj.LinkZonaprop,
			&buildingObj.LinkArgenprop)
		return buildingObj, err
	case "emprendimientos":
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
			pq.Array(&buildingObj.Images),
			&buildingObj.LinkML,
			&buildingObj.LinkZonaprop,
			&buildingObj.LinkArgenprop)
		return buildingObj, err
	default:
		return nil, fmt.Errorf("unsupported category: %s", category)
	}
}

/*
GenerateSQLquery returns "query" and "args"(arguments),
being the $x's inside "query" subsequently replaced by the "args" string values .
*/
func generateSQLquery(category string, urlQyParams map[string][]string) (string, []interface{}) {
	query := fmt.Sprintf(`SELECT * FROM public."%s"`, category)
	args := []interface{}{}
	conditions := []string{}
	var queried []string
	expressionMapping := map[string]string{
		"price_init":            "price >=",
		"price_limit":           "price <=",
		"env_init":              "env >=",
		"env_limit":             "env <=",
		"bedroom_init":          "bedrooms >=",
		"bedroom_limit":         "bedrooms <=",
		"bathroom_init":         "bathrooms >=",
		"bathroom_limit":        "bathrooms <=",
		"garage_init":           "garages >=",
		"garage_limit":          "garages <=",
		"total_surface_init":    "total_surface >=",
		"total_surface_limit":   "total_surface <=",
		"covered_surface_init":  "covered_surface >=",
		"covered_surface_limit": "covered_surface <=",
	}

	for fieldKey, fieldValue := range urlQyParams {

		// ************************************************************************
		// Mgmt. of attributes within expressionMapping (share same query syntax)
		if expression, ok := expressionMapping[fieldKey]; ok {
			conditions = append(
				conditions, fmt.Sprintf("(%s $%d %s)",
					expression,
					len(args)+1,
					func() string {
						wasQueried := false
						words := strings.Fields(expression)
						for _, value := range queried {
							if words[0] == value {
								wasQueried = true
							}
						}
						if wasQueried {
							return ""
						} else {
							queried = append(queried, words[0])
							return fmt.Sprintf("OR %s IS NULL", words[0])
						}
					}()))
			args = append(args, fieldValue[0])

			// ************************************************************************
			// Mgmt. of the rest of the attributes (distinctive query syntax)
		} else if fieldKey == "location" {
			conditions = append(conditions, fmt.Sprintf("location ILIKE $%d", len(args)+1))
			args = append(args, "%"+fieldValue[0]+"%")
		} else if fieldKey == "building_status" {
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
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//********************************
	// Slice of structs initialization.
	var buildings []interface{}
	for rows.Next() {
		newBuilding, err := initBuildingType(category, rows)
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
	categoryHandler := func(w http.ResponseWriter, r *http.Request) {
		getDBdata(w, r, db)
	}
	sv.Get("/{category}", categoryHandler)

	//******************************************
	// Turning on the server.
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

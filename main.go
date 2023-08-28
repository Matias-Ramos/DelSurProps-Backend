package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"github.com/Matias-Ramos/Inmobiliaria-backend-go/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
InitBuilingType mutates a *sql.Rows into a Go interface{}.
Such result will represent the building object.
*/
func initBuildingType(category string, rows *sql.Rows) (interface{}, error) {
	switch category {
	case "alquiler_inmuebles":
		buildingObj := &models.RentBuilding{Building: &models.Building{}}
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
			&buildingObj.LinkArgenprop,
			&buildingObj.Currency)
		return buildingObj, err
	case "venta_inmuebles":
		buildingObj := &models.SalesBuilding{Building: &models.Building{}}
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
		buildingObj := &models.VentureBuilding{Building: &models.Building{}}
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
generateGetQuery returns "query" and "args".
The returned "query" contains placeholders like $1, $2, which will be replaced by the values in "args" respectively.
*/
func generateGetQuery(category string, urlQyParams map[string][]string) (string, []interface{}) {
	query := fmt.Sprintf(`SELECT * FROM public."%s"`, category)
	args := []interface{}{}
	conditions := []string{}
	var queried []string
	// The attributes within expressionMapping share same SQL query syntax so I grouped them up.
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

	// DB data gathering through SQL querying.
	category := chi.URLParam(r, "category")
	urlQyParams := r.URL.Query()
	query, args := generateGetQuery(category, urlQyParams)
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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

	// Convertion from Go slice to JSON.
	jsonData, err := json.Marshal(buildings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sending the data to the requester.
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func generateInsertQuery(data interface{}, tableName string) string {
	var columns []string
	var values []string

	buildingType := reflect.TypeOf(data)
	buildingValue := reflect.ValueOf(data).Elem()

	for i := 0; i < buildingType.Elem().NumField(); i++ {
		field := buildingType.Elem().Field(i)
		value := buildingValue.Field(i).Interface()

		columns = append(columns, field.Name)
		if field.Type == reflect.TypeOf([]string{}) {
			values = append(values, fmt.Sprintf("'{%s}'", strings.Join(value.([]string), ",")))
		} else {
			values = append(values, fmt.Sprintf("'%v'", value))
		}
	}

	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	query := fmt.Sprintf("INSERT INTO public.%s (%s) VALUES (%s);", tableName, columnsStr, valuesStr)
	return query
}

/*
postData processes incoming form data:
- Determines the appropriate struct based on the category
- Adds an ID to the form data and populates the corresponding struct
- Calls the query generator with the resulting struct
*/
func postData(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")

	var buildingObj interface{}
	switch category {
	case "alquiler_inmueble":
		buildingObj = &models.TestBuilding{}
	case "venta_inmueble":
		buildingObj = &models.SalesBuilding{Building: &models.Building{}}
	case "emprendimiento":
		buildingObj = &models.VentureBuilding{Building: &models.Building{}}
	}

	// io.ReadCloser to []byte
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error on io.ReadAll() - post method: %s", err)
		return
	}

	// []byte to map
	var m map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &m); err != nil {
		fmt.Printf("error on json.Unmarshal() - post method: %s", err)
		return
	}

	// Add an Id to the map
	id := int64(time.Now().UnixNano()) + int64(rand.Intn(1000000))
	m["Id"] = id

	// Recreate the JSON
	newData, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error on json.Marshal() - post method: %s", err)
		return
	}

	fmt.Println("newData")
	fmt.Println(string(newData))

	// Populate the buildingObj
	if err := json.Unmarshal(newData, buildingObj); err != nil {
		fmt.Printf("error on json.Unmarshal() - post method: %s", err)
		return
	}

	fmt.Println(generateInsertQuery(buildingObj, "alquiler_inmuebles"))
	w.Write([]byte("Okay"))
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
		getDBdata(w, r, db)
	}
	sv.Get("/api/{category}", categoryHandler)
	sv.Post("/admin/post/{category}", postData)

	//******************************************
	// Turning on the server.
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", sv)
}

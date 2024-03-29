package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/models"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)



func GetDBdata(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// *******************************
		// DB data gathering through SQL querying.
		category := chi.URLParam(r, "category")
		urlQyParams := r.URL.Query()
		query, args := generateGetQuery(category, urlQyParams)
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("err at db.Query - GetDBdata: %v", err)
			return
		}
		defer rows.Close()

		// *******************************
		// Slice of structs initialization.
		var buildings []interface{}
		for rows.Next() {
			newBuilding, err := initBuildingType(category, rows)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("err at rows.Next() -> initBuildingType - GetDBdata: %v", err)
				return
			}
			buildings = append(buildings, newBuilding)
		}

		// *******************************
		// Convertion from Go slice to JSON.
		jsonData, err := json.Marshal(buildings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("err at json.Marshal() - GetDBdata: %v", err)
			return
		}

		// *******************************
		// Sending the data to the requester.
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

// **************************
// The attributes that share the same SQL query syntax are grouped in this map.
var commonFieldsMapping = map[string]string{
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
/*
1. generateGetQuery returns "query" and "args".
The returned "query" contains placeholders like $1, $2, which will be replaced by the values in "args" respectively.
*/
func generateGetQuery(category string, urlQyParams map[string][]string) (string, []interface{}) {
	query := fmt.Sprintf(`SELECT * FROM public."%s"`, category)
	args := []interface{}{}
	sqlConditions := []string{}
	var queriedProperties []string

	for fieldKey, fieldValue := range urlQyParams {
		if _, isCommonField := commonFieldsMapping[fieldKey]; isCommonField {
			args, sqlConditions, queriedProperties = handleCommonField(args, sqlConditions, queriedProperties, fieldKey, fieldValue)
		} else if fieldKey == "location" {
			args, sqlConditions = handleLocationField(args, sqlConditions, fieldValue)
		} else if fieldKey == "building_status" {
			args, sqlConditions = handleBuildingStatusField(args, sqlConditions, fieldValue)
		}
	}
	// Add the WHERE clause if there are sqlConditions
	if len(sqlConditions) > 0 {
		query += " WHERE " + strings.Join(sqlConditions, " AND ")
	}
	return query, args
}

// **************************
// 1.1- Argument generators.

func handleCommonField(args []interface{}, sqlConditions []string, queriedProperties []string, fieldKey string, fieldValue []string) ([]interface{}, []string, []string) {
	expression := commonFieldsMapping[fieldKey]
	sqlConditions = append(
		sqlConditions, fmt.Sprintf("(%s $%d %s)",
			expression,
			len(args)+1,
			includeNullValues(expression, queriedProperties)))
	args = append(args, fieldValue[0])
	return args, sqlConditions, queriedProperties
}
func includeNullValues(expression string, queriedProperties []string) string {
	wasQueried := false
	words := strings.Fields(expression)
	for _, value := range queriedProperties {
		if words[0] == value {
			wasQueried = true
		}
	}
	if wasQueried {
		return ""
	} else {
		return fmt.Sprintf("OR %s IS NULL", words[0])
	}
}
func handleLocationField(args []interface{}, sqlConditions []string, fieldValue []string) ([]interface{}, []string) {
	args = append(args, "%"+fieldValue[0]+"%")
	sqlConditions = append(sqlConditions, fmt.Sprintf("location ILIKE $%d", len(args)))
	return args, sqlConditions
}
func handleBuildingStatusField(args []interface{}, sqlConditions []string, fieldValue []string) ([]interface{}, []string) {
	switch fieldValue[0] {
	case "in_progress":
		args = append(args, true)
		sqlConditions = append(sqlConditions, fmt.Sprintf("in_progress = $%d", len(args)))
	case "pozo":
		args = append(args, true)
		sqlConditions = append(sqlConditions, fmt.Sprintf("pozo = $%d", len(args)))
	case "in_progress-or-pozo":
		args = append(args, true, true)
		sqlConditions = append(sqlConditions, fmt.Sprintf("in_progress = $%d OR pozo = $%d", len(args)-1, len(args)))
	}
	return args, sqlConditions
}

// **************************
/*
2. initBuilingType mutates the *sql.Rows from generateGetQuery() into a Go interface{}.
Such result will hold the building properties.
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
			&buildingObj.Link_ml,
			&buildingObj.Link_zonaprop,
			&buildingObj.Link_argenprop,
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
			&buildingObj.Link_ml,
			&buildingObj.Link_zonaprop,
			&buildingObj.Link_argenprop)
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
			&buildingObj.Link_ml,
			&buildingObj.Link_zonaprop,
			&buildingObj.Link_argenprop)
		return buildingObj, err
	default:
		return nil, fmt.Errorf("unsupported category: %s", category)
	}
}

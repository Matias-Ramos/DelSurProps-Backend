package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/models"

	"github.com/go-chi/chi/v5"
)

/*
postData processes incoming form data:
- Determines the appropriate struct based on the category.
- Adds an ID to the form data and populates the corresponding struct.
- Calls the query generator with the resulting struct.
- Executes the query in the DB.
*/
func PostData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := chi.URLParam(r, "category")

		var buildingObj interface{}
		switch category {
		case "alquiler_inmueble":
			buildingObj = &models.RentBuilding{}
		case "venta_inmueble":
			buildingObj = &models.SalesBuilding{Building: &models.Building{}}
		case "emprendimiento":
			buildingObj = &models.VentureBuilding{Building: &models.Building{}}
		}

		// io.ReadCloser to []byte
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on io.ReadAll() - post method: %s", err)
			return
		}

		// []byte to map
		var m map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on json.Unmarshal() - []byte to map - post method: %s", err)
			return
		}

		// Add an Id to the map
		id := int64(rand.Intn(999999999999999))
		m["id"] = id

		priceInt, _ := strconv.Atoi(m["price"].(string))
		m["price"] = priceInt

		// Recreate the JSON
		newData, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on json.Marshal() - post method: %s", err)
			return
		}

		// Populate the buildingObj
		if err := json.Unmarshal(newData, buildingObj); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on json.Unmarshal() - post method: %s", err)
			return
		}

		query := generateInsertQuery(buildingObj, category)
		_, err = db.Exec(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on db.Exec() - post method: %s", err)
		}
	}
}

/*
Note that in models.go, the Building struct is embedded into RentBuilding, SalesBuilding and VentureBuilding.
"generateInsertQuery" iterates through fields at both top level (RB,SB,VB) and low level (Building) fields.
*/
func generateInsertQuery(data interface{}, tableName string) string {
	var sqlColumns []string
	var sqlValues []string
	building := reflect.ValueOf(data).Elem()

	for i := 0; i < building.NumField(); i++ {
		externalField := building.Type().Field(i)
		externalValue := building.Field(i).Interface()

		// If the field is a pointer to Building...
		if externalField.Type.Kind() == reflect.Ptr && externalField.Type.Elem().Name() == "Building" {
			embeddedValues := reflect.Indirect(reflect.ValueOf(externalValue))
			// ...iterate over its internal struct
			for j := 0; j < embeddedValues.NumField(); j++ {
				internalField := embeddedValues.Type().Field(j)
				internalValue := embeddedValues.Field(j).Interface()
				sqlColumns = append(sqlColumns, internalField.Name)
				sqlValues = append(sqlValues, parseToSqlSyntax(internalField, internalValue))
			}
		} else {
			// top-level fields on RentBuilding / SalesBuilding / VentureBuilding
			sqlColumns = append(sqlColumns, externalField.Name)
			sqlValues = append(sqlValues, parseToSqlSyntax(externalField, externalValue))
		}
	}

	sqlColumnsFormatted := strings.Join(convertToLowerCase(sqlColumns), ", ")
	sqlValuesFormatted := strings.Join(sqlValues, ", ")

	query := fmt.Sprintf("INSERT INTO public.%ss (%s) VALUES (%s);", tableName, sqlColumnsFormatted, sqlValuesFormatted)
	return query
}
func convertToLowerCase(sqlColumns []string) []string {
	var lowercaseColumns []string
	for _, col := range sqlColumns {
		lowercaseColumns = append(lowercaseColumns, strings.ToLower(col))
	}
	return lowercaseColumns
}
func parseToSqlSyntax(internalField reflect.StructField, internalValue interface{}) string {
	if isImgsField := internalField.Type == reflect.TypeOf([]string{}); isImgsField {
		return fmt.Sprintf("'{%s}'", strings.Join(internalValue.([]string), ","))
	}

	switch v := internalValue.(type) {
	case models.NullInt16:
		if v.Valid {
			return fmt.Sprintf("'%v'", v.Int16)
		}
		return "NULL"
	case models.NullString:
		if v.Valid {
			return fmt.Sprintf("'%v'", v.String)
		}
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

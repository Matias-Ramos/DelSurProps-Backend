package crud

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Matias-Ramos/Inmobiliaria-backend-go/models"
	
	"github.com/go-chi/chi/v5"
)


/*
postData processes incoming form data:
- Determines the appropriate struct based on the category
- Adds an ID to the form data and populates the corresponding struct
- Calls the query generator with the resulting struct
*/
func PostData(w http.ResponseWriter, r *http.Request) {
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

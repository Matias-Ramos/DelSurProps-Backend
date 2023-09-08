package crud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func DeleteData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// io.ReadCloser to []byte
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on io.ReadAll() - delete method: %s", err)
			return
		}

		// []byte to map
		var m map[string]string
		if err := json.Unmarshal(bodyBytes, &m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on json.Unmarshal() - []byte to map - delete method: %s", err)
			return
		}

		// Generate the query
		buildingId := m["buildingId"]
		category := chi.URLParam(r, "category")
		unslashedCategory := strings.Trim(category, "/")
		query := fmt.Sprintf("DELETE FROM %s WHERE id=%s", unslashedCategory, buildingId)

		// Execute the query
		_, err = db.Exec(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on db.Exec() - delete method: %s", err)
		}
	}
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func generateSQLquery(category string, urlQyParams map[string][]string) (string, []interface{}) {
	// Building the SQL query
	// (this way to query prevents SQL injection vulnerabilities)
	query := fmt.Sprintf(`SELECT * FROM public."%s"`, category)
	args := []interface{}{}
	conditions := []string{}
	columnMapping := map[string]string{
		"location":          "location ILIKE",
		"price_init":        "price >=",
		"price_limit":       "price <=",
		"env_init":          "env >=",
		"env_limit":         "env <=",
		"bedroom_init":      "bedrooms >=",
		"bedroom_limit":     "bedrooms <=",
		"bathroom_init":     "bathrooms >=",
		"bathroom_limit":    "bathrooms <=",
		"garage_init":       "garages >=",
		"garage_limit":      "garages <=",
		"total_surface_init":   "total_surface >=",
		"total_surface_limit":  "total_surface <=",
		"covered_surface_init": "covered_surface >=",
		"covered_surface_limit": "covered_surface <=",
	}

	for fieldKey, fieldValue := range urlQyParams {
		if column, ok := columnMapping[fieldKey]; ok {
			// The field key exists in the mapping, so add the condition
			conditions = append(conditions, fmt.Sprintf("%s $%d", column, len(args)+1))
			args = append(args, fieldValue[0])
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

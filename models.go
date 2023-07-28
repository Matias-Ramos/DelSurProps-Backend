package main

import "database/sql"

type Building struct {
	Id            int            `json:"id"`
	Location      string         `json:"location"`
	Price         int            `json:"price"`
	Env           sql.NullInt16  `json:"env"`
	Bedrooms      sql.NullInt16  `json:"bedrooms"`
	Bathrooms     sql.NullInt16  `json:"bathrooms"`
	Garages       sql.NullInt16  `json:"garages"`
	LinkML        sql.NullString `json:"linkML"`
	LinkZonaprop  sql.NullString `json:"linkZonaprop"`
	LinkArgenprop sql.NullString `json:"linkArgenprop"`
	Images        []string       `json:"images"`
}
type RentBuilding struct {
	*Building
}
type SalesBuilding struct {
	*Building
	Covered_surface sql.NullInt16 `json:"covered_surface"`
	Total_surface   sql.NullInt16 `json:"total_surface"`
}
type VentureBuilding struct {
	*Building
	Covered_surface sql.NullInt16 `json:"covered_surface"`
	Total_surface   sql.NullInt16 `json:"total_surface"`
	Pozo            bool          `json:"pozo"`
	In_progress     bool          `json:"in_progress"`
}

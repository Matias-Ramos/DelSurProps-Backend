package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// **************************************
// custom React null receiving data types.
type NullInt16 struct {
	sql.NullInt16
}
type NullString struct {
	sql.NullString
}

// ******************************************************
// UnmarshalJSON implements the json.Unmarshaler interface

func (niP *NullInt16) UnmarshalJSON(data []byte) error {
	stringedValue := strings.Trim(string(data), `"`)
	if len(stringedValue) == 0 {
		niP.Valid = false
		return nil
	}

	// byte string -> string -> int
	intValue, errIntConv := strconv.Atoi(stringedValue)
	if errIntConv != nil {
		return fmt.Errorf("UnmarshalJSON errIntConv at models.go failed: %v", errIntConv)
	}
	niP.Int16 = int16(intValue)
	niP.Valid = true

	return nil
}

func (nsP *NullString) UnmarshalJSON(data []byte) error {
	stringedValue := strings.Trim(string(data), `"`)
	if len(stringedValue) == 0 {
		nsP.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &nsP.String)
	if err == nil {
		nsP.Valid = true
		return nil
	}
	return err
}

// **********************
// Building structures

type TestBuilding struct {
	Id            int        `json:"id"`
	Location      string     `json:"location"`
	Price         NullInt16  `json:"price"`
	Env           NullInt16  `json:"env"`
	Bedrooms      NullInt16  `json:"bedrooms"`
	Bathrooms     NullInt16  `json:"bathrooms"`
	Garages       NullInt16  `json:"garages"`
	LinkML        NullString `json:"link_ml"`
	LinkZonaprop  NullString `json:"link_zonaprop"`
	LinkArgenprop NullString `json:"link_argenprop"`
	Images        []string   `json:"image_links"`
	Currency      string     `json:"currency"`
}
type TestRentBuilding struct {
	*TestBuilding
	Currency string `json:"currency"`
}
type TestSalesBuilding struct {
	*TestBuilding
	Covered_surface sql.NullInt16 `json:"covered_surface"`
	Total_surface   sql.NullInt16 `json:"total_surface"`
}
type TestVentureBuilding struct {
	*TestBuilding
	Covered_surface sql.NullInt16 `json:"covered_surface"`
	Total_surface   sql.NullInt16 `json:"total_surface"`
	Pozo            bool          `json:"pozo"`
	In_progress     bool          `json:"in_progress"`
}

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
	Currency string `json:"currency"`
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

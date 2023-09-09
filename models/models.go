package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// **************************************
// Custom null receiving data types.
type NullInt16 struct {
	sql.NullInt16
}
type NullString struct {
	sql.NullString
}

// ******************************************************
// UnmarshalJSON implements the json.Unmarshal interface

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

type Building struct {
	Id             int64      `json:"id"`
	Location       string     `json:"location"`
	Price          int        `json:"price"`	
	Env            NullInt16  `json:"env"`
	Bedrooms       NullInt16  `json:"bedrooms"`
	Bathrooms      NullInt16  `json:"bathrooms"`
	Garages        NullInt16  `json:"garages"`
	Link_ml        NullString `json:"link_ml"`
	Link_zonaprop  NullString `json:"link_zonaprop"`
	Link_argenprop NullString `json:"link_argenprop"`
	Images         []string   `json:"image_links"`
}
type RentBuilding struct {
	*Building
	Currency string `json:"currency"`
}
type SalesBuilding struct {
	*Building
	Covered_surface NullInt16 `json:"covered_surface"`
	Total_surface   NullInt16 `json:"total_surface"`
}
type VentureBuilding struct {
	*Building
	Covered_surface NullInt16 `json:"covered_surface"`
	Total_surface   NullInt16 `json:"total_surface"`
	Pozo            bool      `json:"pozo"`
	In_progress     bool      `json:"in_progress"`
}

package main

type Building struct {
	Id        int      `json:"id"`
	Location  string   `json:"location"`
	Price     int      `json:"price"`
	Env      uint8    `json:"env"`
	Bedrooms  uint8    `json:"bedrooms"`
	Bathrooms uint8    `json:"bathrooms"`
	Garages   uint8    `json:"garages"`
	Images    []string `json:"images"`
}
type RentBuilding struct {
	*Building
}
type SalesBuilding struct {
	*Building
	Covered_surface int `json:"covered_surface"`
	Total_surface   int `json:"total_surface"`
}
type VentureBuilding struct {
	*Building
	Covered_surface int  `json:"covered_surface"`
	Total_surface   int  `json:"total_surface"`
	Pozo            bool `json:"pozo"`
	In_progress     bool `json:"in_progress"`
}

package main

type Building struct {
	Id          int      `json:"id"`
	Ubicacion   string   `json:"ubicacion"`
	Precio      int      `json:"precio"`
	Ambientes   uint8    `json:"ambientes"`
	Dormitorios uint8    `json:"dormitorios"`
	Banios      uint8    `json:"banios"`
	Garages     uint8    `json:"garages"`
	Imagenes    []string `json:"imagenes"`
}
type RentBuilding struct {
	*Building
}
type SalesBuilding struct {
	*Building
	Superficie_cubierta int `json:"superficie_cubierta"`
	Superficie_total    int `json:"superficie_total"`
}
type VentureBuilding struct {
	*Building
	Superficie_cubierta int  `json:"superficie_cubierta"`
	Superficie_total    int  `json:"superficie_total"`
	En_pozo             bool `json:"en_pozo"`
	En_construccion     bool `json:"en_construccion"`
}

package venta

import (
	"encoding/json"
)

type Venta struct {
	Cliente  string    `json:"cliente"`
	Cantidad string    `json:"cantidad"`
	Hora     string    `json:"hora"`
	Stock    string    `json:"stock"`
	Coords   []float64 `json:"ubicacion"`
}

func (v *Venta) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

package venta

import (
	"encoding/json"
)

type Venta struct {
	Maestro  string     `json:"maestro,omitempty"`
	Cliente  string     `json:"cliente"`
	Cantidad int        `json:"cantidad"`
	Hora     string     `json:"hora"`
	Stock    int        `json:"stock"`
	Coords   [2]float64 `json:"coords"`
}

func (v *Venta) JSON() ([]byte, error) {
	return json.Marshal(v)
}

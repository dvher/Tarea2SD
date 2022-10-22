package miembro

import (
	"encoding/json"
)

type Miembro struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Rut      string `json:"rut"`
	Email    string `json:"email"`
	Patente  string `json:"patente"`
	Premium  bool   `json:"premium"`
}

func (m *Miembro) JSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Miembro) JSONIndent() ([]byte, error) {
	return json.MarshalIndent(m, "", "\t")
}

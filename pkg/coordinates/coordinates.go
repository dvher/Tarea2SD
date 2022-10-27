package coordinates

import (
	"encoding/json"
)

type Coordinates struct {
	Coords  [2]float64 `json:"coords"`
	Miembro string     `json:"rut,omitempty"`
}

func (c *Coordinates) JSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Coordinates) JSONIndent() ([]byte, error) {
	return json.MarshalIndent(c, "", "\t")
}

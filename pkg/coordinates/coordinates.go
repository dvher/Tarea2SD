package coordinates

import (
	"encoding/json"
)

type Coordinates struct {
	Coords [2]float64 `json:"coords"`
}

func (c *Coordinates) JSON() ([]byte, error) {
	return json.Marshal(c)
}

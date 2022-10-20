package coordinates

import (
    "encoding/json"
)

type Coordinates struct {
   Coords []float64 `json:"coords"`
}

func (c *Coordinates) MarshalJSON() ([]byte, error) {
    return json.Marshal(c)
}

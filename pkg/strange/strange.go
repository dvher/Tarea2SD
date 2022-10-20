package strange

import (
    "encoding/json"
)

type Strange struct {
   Coords []float64 `json:"coords"`
}

func (s *Strange) MarshalJSON() ([]byte, error) {
    return json.Marshal(s)
}

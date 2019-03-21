package encode

import (
	"encoding/json"
	"net/http"
)

type jsonEncoder struct{}

//NewJSONEncoder for JSON API endpoints
func NewJSONEncoder() Encoder {
	return &jsonEncoder{}
}

func (e *jsonEncoder) Encode(w http.ResponseWriter, v interface{}) error {
	defer w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}

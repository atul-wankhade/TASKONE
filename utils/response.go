package utils

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

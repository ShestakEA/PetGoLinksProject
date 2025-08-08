package response

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, res any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}

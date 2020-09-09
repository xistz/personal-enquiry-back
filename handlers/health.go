package handlers

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles requests to "/health"
func HealthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	health := healthResponse{
		Status: "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&health)
}

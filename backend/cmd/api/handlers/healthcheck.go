package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck handles GET /v1/healthcheck requests.
// Returns 200 OK with {"status": "ok"} when healthy.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

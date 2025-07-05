// handlers/system.go
// Contains handlers for system-level checks, like health.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// HealthCheckHandler checks the health of the service, including the database connection.
func (env *Env) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// The Ping method was added to the Store interface for this purpose.
	if err := env.DB.Ping(r.Context()); err != nil {
		log.Error().Err(err).Msg("Health check failed: database is unreachable")
		// If pinging the DB fails, return a 503 Service Unavailable
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "down",
			"error":  "database is unreachable",
		})
		return
	}

	// Everything is okay
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "up"})
}

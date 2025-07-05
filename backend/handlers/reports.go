// handlers/reports.go
// Contains handlers for generating financial reports.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"your_username/budget-tracker/auth"
)

func (env *Env) GetMonthlySummaryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get year and month from query params, with defaults
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year == 0 {
		year = time.Now().Year()
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month == 0 {
		month = int(time.Now().Month())
	}

	summary, err := env.DB.GetMonthlySummary(userID, year, month)
	if err != nil {
		http.Error(w, "Failed to generate summary report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

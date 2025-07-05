// handlers/budgets.go
// Contains handlers for managing monthly budgets.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/database"
	"your_username/budget-tracker/models"
	"your_username/budget-tracker/utils"

	"github.com/gorilla/mux"
)

func (env *Env) CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var budget models.Budget
	if err := utils.ParseAndValidate(r, &budget); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	budget.UserID = userID

	if err := env.DB.CreateBudget(&budget); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			http.Error(w, "A budget for this category and month already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create budget", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budget)
}

func (env *Env) GetBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	budgets, err := env.DB.GetBudgets(userID, year, month)
	if err != nil {
		http.Error(w, "Failed to retrieve budgets", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func (env *Env) UpdateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	var budget models.Budget
	if err := utils.ParseAndValidate(r, &budget); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	budget.ID = id
	budget.UserID = userID

	if err := env.DB.UpdateBudget(&budget); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Budget not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update budget", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budget)
}

func (env *Env) DeleteBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	if err := env.DB.DeleteBudget(id, userID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Budget not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete budget", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

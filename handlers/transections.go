// handlers/transactions.go
// Contains handlers for transaction-related operations, including filtering and sorting.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/database"
	"your_username/budget-tracker/models"
	"your_username/budget-tracker/utils"

	"github.com/gorilla/mux"
)

func (env *Env) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var t models.Transaction
	if err := utils.ParseAndValidate(r, &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.UserID = userID

	if err := env.DB.CreateTransaction(&t); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (env *Env) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r, env.DefaultPageLimit)
	filter := parseTransactionFilters(r)
	filter.Sort = parseSort(r.URL.Query().Get("sort"), env.AllowedSortFields)

	transactions, totalRecords, err := env.DB.GetTransactions(userID, limit, offset, filter)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	metadata := utils.CalculateMetadata(totalRecords, page, limit)
	response := models.PaginatedResponse{
		Metadata: metadata,
		Data:     transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (env *Env) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	transaction, err := env.DB.GetTransactionByID(id, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Transaction not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve transaction", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (env *Env) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	var t models.Transaction
	if err := utils.ParseAndValidate(r, &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := env.DB.UpdateTransaction(id, userID, &t); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Transaction not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		}
		return
	}
	t.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (env *Env) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	if err := env.DB.DeleteTransaction(id, userID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Transaction not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseTransactionFilters(r *http.Request) models.TransactionFilter {
	q := r.URL.Query()
	filter := models.TransactionFilter{
		Type: q.Get("type"),
	}
	if min, err := strconv.ParseFloat(q.Get("min_amount"), 64); err == nil {
		filter.MinAmount = &min
	}
	if max, err := strconv.ParseFloat(q.Get("max_amount"), 64); err == nil {
		filter.MaxAmount = &max
	}
	if start, err := time.Parse(time.RFC3339, q.Get("start_date")); err == nil {
		filter.StartDate = &start
	}
	if end, err := time.Parse(time.RFC3339, q.Get("end_date")); err == nil {
		filter.EndDate = &end
	}
	return filter
}

func parseSort(sortQuery string, allowedFields []string) string {
	if sortQuery == "" {
		return "ORDER BY date DESC"
	}
	parts := strings.Split(sortQuery, ",")
	var orderByClauses []string
	for _, part := range parts {
		direction := "ASC"
		field := part
		if strings.HasPrefix(part, "-") {
			direction = "DESC"
			field = part[1:]
		}
		isAllowed := false
		for _, allowed := range allowedFields {
			if field == allowed {
				isAllowed = true
				break
			}
		}
		if isAllowed {
			orderByClauses = append(orderByClauses, field+" "+direction)
		}
	}
	if len(orderByClauses) == 0 {
		return "ORDER BY date DESC"
	}
	return "ORDER BY " + strings.Join(orderByClauses, ", ")
}

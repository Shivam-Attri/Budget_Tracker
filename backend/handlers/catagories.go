// handlers/categories.go
// Contains handlers for managing custom categories.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/database"
	"your_username/budget-tracker/models"
	"your_username/budget-tracker/utils"

	"github.com/gorilla/mux"
)

func (env *Env) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var cat models.Category
	if err := utils.ParseAndValidate(r, &cat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cat.UserID = userID

	if err := env.DB.CreateCategory(&cat); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			http.Error(w, "Category with this name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (env *Env) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r, env.DefaultPageLimit)

	categories, totalRecords, err := env.DB.GetCategories(userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}

	metadata := utils.CalculateMetadata(totalRecords, page, limit)
	response := models.PaginatedResponse{
		Metadata: metadata,
		Data:     categories,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (env *Env) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	var cat models.Category
	if err := utils.ParseAndValidate(r, &cat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cat.ID = id
	cat.UserID = userID

	if err := env.DB.UpdateCategory(&cat); err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			http.Error(w, "Category with this name already exists", http.StatusConflict)
		} else if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cat)
}

func (env *Env) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	id := params["id"]

	if err := env.DB.DeleteCategory(id, userID); err != nil {
		if errors.Is(err, database.ErrCategoryInUse) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

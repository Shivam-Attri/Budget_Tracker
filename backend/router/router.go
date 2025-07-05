// router/router.go
// Defines all API endpoints and applies middleware.
package router

import (
	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/handlers"
	"your_username/budget-tracker/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(env *handlers.Env) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware first.
	r.Use(middleware.RateLimiter)

	// Public routes
	r.HandleFunc("/health", env.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/register", env.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", env.LoginHandler).Methods("POST")
	r.HandleFunc("/refresh", env.RefreshTokenHandler).Methods("POST")

	// Protected API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(auth.Middleware) // Auth middleware is specific to this subrouter

	api.HandleFunc("/logout", env.LogoutHandler).Methods("POST")

	// Categories Routes
	api.HandleFunc("/categories", env.CreateCategoryHandler).Methods("POST")
	api.HandleFunc("/categories", env.GetCategoriesHandler).Methods("GET")
	api.HandleFunc("/categories/{id}", env.UpdateCategoryHandler).Methods("PUT")
	api.HandleFunc("/categories/{id}", env.DeleteCategoryHandler).Methods("DELETE")

	// Transactions Routes
	api.HandleFunc("/transactions", env.CreateTransactionHandler).Methods("POST")
	api.HandleFunc("/transactions", env.GetTransactionsHandler).Methods("GET")
	api.HandleFunc("/transactions/{id}", env.GetTransactionByIDHandler).Methods("GET")
	api.HandleFunc("/transactions/{id}", env.UpdateTransactionHandler).Methods("PUT")
	api.HandleFunc("/transactions/{id}", env.DeleteTransactionHandler).Methods("DELETE")

	// Budgets Routes
	api.HandleFunc("/budgets", env.CreateBudgetHandler).Methods("POST")
	api.HandleFunc("/budgets", env.GetBudgetsHandler).Methods("GET")
	api.HandleFunc("/budgets/{id}", env.UpdateBudgetHandler).Methods("PUT")
	api.HandleFunc("/budgets/{id}", env.DeleteBudgetHandler).Methods("DELETE")

	// Reports Routes
	api.HandleFunc("/reports/summary", env.GetMonthlySummaryHandler).Methods("GET")

	return r
}

// database/store.go
// Implements the Store interface for all database operations.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/encryption"
	"your_username/budget-tracker/models"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrCategoryInUse = errors.New("category is in use by a transaction")
	ErrDuplicate     = errors.New("duplicate resource")
)

// Store is the complete interface for all database operations.
type Store interface {
	Ping(ctx context.Context) error
	// User methods
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	// Refresh Token methods
	StoreRefreshToken(userID, token string, expiresAt time.Time) error
	ValidateRefreshToken(userID, token string) (bool, error)
	RevokeRefreshToken(token string) error
	// Category methods
	CreateCategory(cat *models.Category) error
	GetCategories(userID string, limit, offset int) ([]*models.Category, int, error)
	UpdateCategory(cat *models.Category) error
	DeleteCategory(id, userID string) error
	// Transaction methods
	CreateTransaction(t *models.Transaction) error
	GetTransactions(userID string, limit, offset int, filter models.TransactionFilter) ([]*models.Transaction, int, error)
	GetTransactionByID(id, userID string) (*models.Transaction, error)
	UpdateTransaction(id, userID string, t *models.Transaction) error
	DeleteTransaction(id, userID string) error
	// Budget methods
	CreateBudget(budget *models.Budget) error
	GetBudgets(userID string, year, month int) ([]*models.Budget, error)
	UpdateBudget(budget *models.Budget) error
	DeleteBudget(id, userID string) error
	// Report methods
	GetMonthlySummary(userID string, year, month int) (*models.ReportSummary, error)
}

type DBStore struct {
	DB  *sql.DB
	Key []byte
}

func (store *DBStore) Ping(ctx context.Context) error { return store.DB.PingContext(ctx) }

// --- User Methods ---
// CreateUser is now wrapped in a transaction to ensure atomicity.
func (store *DBStore) CreateUser(user *models.User) error {
	tx, err := store.DB.Begin()
	if err != nil {
		return err
	}
	// Defer a rollback in case of a panic. If Commit() succeeds, the rollback is a no-op.
	defer tx.Rollback()

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userQuery := `INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(userQuery, user.ID, user.Email, hashedPassword, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return ErrEmailExists
		}
		return err
	}

	// Create default categories within the same transaction.
	defaultCategories := []string{"Groceries", "Salary", "Rent", "Utilities", "Transport"}
	categoryQuery := `INSERT INTO categories (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`

	for _, catName := range defaultCategories {
		encryptedName, err := encryption.Encrypt([]byte(catName), store.Key)
		if err != nil {
			// If encryption fails, the whole transaction should fail.
			return fmt.Errorf("failed to encrypt default category name: %w", err)
		}
		_, err = tx.Exec(categoryQuery, uuid.New(), user.ID, encryptedName, time.Now())
		if err != nil {
			// If a category insert fails, the whole transaction should fail.
			return fmt.Errorf("failed to insert default category: %w", err)
		}
	}

	// If all operations were successful, commit the transaction.
	return tx.Commit()
}

func (store *DBStore) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := store.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return user, err
}

// --- Refresh Token Methods ---
func (store *DBStore) StoreRefreshToken(userID, token string, expiresAt time.Time) error {
	hash := auth.HashToken(token)
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := store.DB.Exec(query, uuid.New(), userID, hash, expiresAt)
	return err
}

func (store *DBStore) ValidateRefreshToken(userID, token string) (bool, error) {
	hash := auth.HashToken(token)
	var storedUserID string
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`
	err := store.DB.QueryRow(query, hash).Scan(&storedUserID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if storedUserID != userID {
		return false, nil
	}
	if time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}

func (store *DBStore) RevokeRefreshToken(token string) error {
	hash := auth.HashToken(token)
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := store.DB.Exec(query, hash)
	return err
}

// --- Category Methods ---
func (store *DBStore) CreateCategory(cat *models.Category) error {
	cat.ID = uuid.New().String()
	cat.CreatedAt = time.Now()
	encryptedName, err := encryption.Encrypt([]byte(cat.Name), store.Key)
	if err != nil {
		return err
	}
	query := `INSERT INTO categories (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err = store.DB.Exec(query, cat.ID, cat.UserID, encryptedName, cat.CreatedAt)
	if err != nil && strings.Contains(err.Error(), "categories_user_id_name_key") {
		return ErrDuplicate
	}
	return err
}

func (store *DBStore) GetCategories(userID string, limit, offset int) ([]*models.Category, int, error) {
	var totalRecords int
	countQuery := "SELECT COUNT(*) FROM categories WHERE user_id = $1"
	if err := store.DB.QueryRow(countQuery, userID).Scan(&totalRecords); err != nil {
		return nil, 0, err
	}
	if totalRecords == 0 {
		return []*models.Category{}, 0, nil
	}
	query := "SELECT id, name, created_at FROM categories WHERE user_id = $1 ORDER BY name LIMIT $2 OFFSET $3"
	rows, err := store.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var categories []*models.Category
	for rows.Next() {
		cat := &models.Category{UserID: userID}
		var encryptedName string
		if err := rows.Scan(&cat.ID, &encryptedName, &cat.CreatedAt); err != nil {
			return nil, 0, err
		}
		decrypted, err := encryption.Decrypt(encryptedName, store.Key)
		if err != nil {
			return nil, 0, err
		}
		cat.Name = string(decrypted)
		categories = append(categories, cat)
	}
	return categories, totalRecords, nil
}

func (store *DBStore) UpdateCategory(cat *models.Category) error {
	encryptedName, err := encryption.Encrypt([]byte(cat.Name), store.Key)
	if err != nil {
		return err
	}
	query := `UPDATE categories SET name = $1 WHERE id = $2 AND user_id = $3`
	res, err := store.DB.Exec(query, encryptedName, cat.ID, cat.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "categories_user_id_name_key") {
			return ErrDuplicate
		}
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (store *DBStore) DeleteCategory(id, userID string) error {
	query := "DELETE FROM categories WHERE id = $1 AND user_id = $2"
	res, err := store.DB.Exec(query, id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "transactions_category_id_fkey") {
			return ErrCategoryInUse
		}
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Transaction Methods ---
func (store *DBStore) CreateTransaction(t *models.Transaction) error {
	t.ID = uuid.New().String()
	encryptedDesc, err := encryption.Encrypt([]byte(t.Description), store.Key)
	if err != nil {
		return err
	}
	query := `INSERT INTO transactions (id, user_id, category_id, description, amount, date, type) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = store.DB.Exec(query, t.ID, t.UserID, t.CategoryID, encryptedDesc, t.Amount, t.Date, t.Type)
	return err
}

func (store *DBStore) GetTransactions(userID string, limit, offset int, filter models.TransactionFilter) ([]*models.Transaction, int, error) {
	var args []interface{}
	var conditions []string
	paramCount := 1
	args = append(args, userID)
	conditions = append(conditions, fmt.Sprintf("user_id = $%d", paramCount))
	paramCount++
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", paramCount))
		args = append(args, filter.Type)
		paramCount++
	}
	if filter.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount >= $%d", paramCount))
		args = append(args, *filter.MinAmount)
		paramCount++
	}
	if filter.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount <= $%d", paramCount))
		args = append(args, *filter.MaxAmount)
		paramCount++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", paramCount))
		args = append(args, *filter.StartDate)
		paramCount++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", paramCount))
		args = append(args, *filter.EndDate)
		paramCount++
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	var totalRecords int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions %s", whereClause)
	if err := store.DB.QueryRow(countQuery, args...).Scan(&totalRecords); err != nil {
		return nil, 0, err
	}
	if totalRecords == 0 {
		return []*models.Transaction{}, 0, nil
	}
	var queryBuilder strings.Builder
	queryBuilder.WriteString(fmt.Sprintf("SELECT id, user_id, category_id, description, amount, date, type FROM transactions %s ", whereClause))
	queryBuilder.WriteString(filter.Sort)
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1))
	args = append(args, limit, offset)
	rows, err := store.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var transactions []*models.Transaction
	for rows.Next() {
		t := &models.Transaction{}
		var encryptedDesc string
		if err := rows.Scan(&t.ID, &t.UserID, &t.CategoryID, &encryptedDesc, &t.Amount, &t.Date, &t.Type); err != nil {
			return nil, 0, err
		}
		decrypted, err := encryption.Decrypt(encryptedDesc, store.Key)
		if err != nil {
			return nil, 0, err
		}
		t.Description = string(decrypted)
		transactions = append(transactions, t)
	}
	return transactions, totalRecords, nil
}

func (store *DBStore) GetTransactionByID(id, userID string) (*models.Transaction, error) {
	t := &models.Transaction{}
	var encryptedDesc string
	query := "SELECT id, user_id, category_id, description, amount, date, type FROM transactions WHERE id = $1 AND user_id = $2"
	err := store.DB.QueryRow(query, id, userID).Scan(&t.ID, &t.UserID, &t.CategoryID, &encryptedDesc, &t.Amount, &t.Date, &t.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	decrypted, err := encryption.Decrypt(encryptedDesc, store.Key)
	if err != nil {
		return nil, err
	}
	t.Description = string(decrypted)
	return t, nil
}

func (store *DBStore) UpdateTransaction(id, userID string, t *models.Transaction) error {
	encryptedDesc, err := encryption.Encrypt([]byte(t.Description), store.Key)
	if err != nil {
		return err
	}
	query := `UPDATE transactions SET category_id = $3, description = $4, amount = $5, date = $6, type = $7 WHERE id = $1 AND user_id = $2`
	res, err := store.DB.Exec(query, id, userID, t.CategoryID, encryptedDesc, t.Amount, t.Date, t.Type)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (store *DBStore) DeleteTransaction(id, userID string) error {
	query := "DELETE FROM transactions WHERE id = $1 AND user_id = $2"
	res, err := store.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Budget Methods ---
func (store *DBStore) CreateBudget(budget *models.Budget) error {
	budget.ID = uuid.New().String()
	budget.CreatedAt = time.Now()
	query := `INSERT INTO budgets (id, user_id, category_id, amount, month, year, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := store.DB.Exec(query, budget.ID, budget.UserID, budget.CategoryID, budget.Amount, budget.Month, budget.Year, budget.CreatedAt)
	if err != nil && strings.Contains(err.Error(), "budgets_user_id_category_id_month_year_key") {
		return ErrDuplicate
	}
	return err
}

func (store *DBStore) GetBudgets(userID string, year, month int) ([]*models.Budget, error) {
	query := `SELECT id, user_id, category_id, amount, month, year, created_at FROM budgets WHERE user_id = $1 AND year = $2 AND month = $3`
	rows, err := store.DB.Query(query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var budgets []*models.Budget
	for rows.Next() {
		b := &models.Budget{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Month, &b.Year, &b.CreatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func (store *DBStore) UpdateBudget(budget *models.Budget) error {
	query := `UPDATE budgets SET amount = $1 WHERE id = $2 AND user_id = $3`
	res, err := store.DB.Exec(query, budget.Amount, budget.ID, budget.UserID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (store *DBStore) DeleteBudget(id, userID string) error {
	query := "DELETE FROM budgets WHERE id = $1 AND user_id = $2"
	res, err := store.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Report Methods ---
func (store *DBStore) GetMonthlySummary(userID string, year, month int) (*models.ReportSummary, error) {
	summary := &models.ReportSummary{}
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
            COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense
        FROM transactions
        WHERE user_id = $1 AND EXTRACT(YEAR FROM date) = $2 AND EXTRACT(MONTH FROM date) = $3;
    `
	err := store.DB.QueryRow(query, userID, year, month).Scan(&summary.TotalIncome, &summary.TotalExpense)
	if err != nil {
		return nil, err
	}
	summary.NetSavings = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

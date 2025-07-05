// models/report.go
// Defines the data structure for the financial summary report.
package models

type ReportSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetSavings   float64 `json:"net_savings"`
}

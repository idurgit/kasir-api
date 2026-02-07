package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTodaySalesSummary() (*models.DailySalesSummary, error) {
	return s.repo.GetTodaySalesSummary()
}

func (s *TransactionService) GetSalesSummaryByDateRange(startDate, endDate string) (*models.DailySalesSummary, error) {
	return s.repo.GetSalesSummaryByDateRange(startDate, endDate)
}

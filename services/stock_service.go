package services

import (
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type StockService struct {
	repo *repositories.StockRepository
}

func NewStockService(repo *repositories.StockRepository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) GetStock(id int) (*db.Stock, error) {
	return s.repo.GetStock(id)
}

func (s *StockService) GetAllStocks() ([]*db.Stock, error) {
	return s.repo.GetAllStocks()
}

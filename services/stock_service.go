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

func (s *StockService) CreateStock(stock *db.Stock) (*db.Stock, error) {
	if err := stock.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateStock(stock)
}

func (s *StockService) UpdateStock(id int, stock *db.Stock) (*db.Stock, error) {
	// Verify stock exists
	_, err := s.repo.GetStock(id)
	if err != nil {
		return nil, err
	}

	if err := stock.Validate(); err != nil {
		return nil, err
	}

	return s.repo.UpdateStock(id, stock)
}

func (s *StockService) DeleteStock(id int) error {
	return s.repo.DeleteStock(id)
}

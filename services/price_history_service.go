package services

import (
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type PriceHistoryService struct {
	repo *repositories.PriceHistoryRepository
}

func NewPriceHistoryService(repo *repositories.PriceHistoryRepository) *PriceHistoryService {
	return &PriceHistoryService{repo: repo}
}

func (s *PriceHistoryService) GetPriceHistoryByStockID(stockID int) ([]*db.PriceHistory, error) {
	return s.repo.GetPriceHistoryByStockID(stockID)
}

func (s *PriceHistoryService) GetPriceHistoryByDateRange(stockID int, startDate, endDate string) ([]*db.PriceHistory, error) {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Set end time to end of day
	endTime = endTime.Add(time.Hour * 24)

	return s.repo.GetPriceHistoryByStockIDAndDateRange(stockID, startTime, endTime)
}

func (s *PriceHistoryService) RecordPrice(stockID int, price float64, date time.Time) (*db.PriceHistory, error) {
	priceHistory := &db.PriceHistory{
		StockID: stockID,
		Price:   price,
		Date:    date,
	}

	if err := priceHistory.Validate(); err != nil {
		return nil, err
	}

	return s.repo.RecordPrice(priceHistory)
}

func (s *PriceHistoryService) GetLatestPrice(stockID int) (*db.PriceHistory, error) {
	return s.repo.GetLatestPrice(stockID)
}

func (s *PriceHistoryService) CalculatePriceStats(stockID int) (map[string]float64, error) {
	priceHistories, err := s.repo.GetPriceHistoryByStockID(stockID)
	if err != nil {
		return nil, err
	}

	if len(priceHistories) == 0 {
		return nil, db.ErrRecordNotFound
	}

	var minPrice, maxPrice, totalPrice float64
	minPrice = priceHistories[0].Price
	maxPrice = priceHistories[0].Price

	for _, ph := range priceHistories {
		if ph.Price < minPrice {
			minPrice = ph.Price
		}
		if ph.Price > maxPrice {
			maxPrice = ph.Price
		}
		totalPrice += ph.Price
	}

	avgPrice := totalPrice / float64(len(priceHistories))

	stats := map[string]float64{
		"min":     minPrice,
		"max":     maxPrice,
		"average": avgPrice,
		"latest":  priceHistories[0].Price,
	}

	return stats, nil
}

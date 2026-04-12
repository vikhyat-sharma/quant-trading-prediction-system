package services

import (
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type PortfolioService struct {
	repo *repositories.PortfolioRepository
}

func NewPortfolioService(repo *repositories.PortfolioRepository) *PortfolioService {
	return &PortfolioService{repo: repo}
}

func (s *PortfolioService) GetPortfoliosByUserID(userID int) ([]*db.Portfolio, error) {
	return s.repo.GetPortfoliosByUserID(userID)
}

func (s *PortfolioService) GetPortfolioByID(userID, portfolioID int) (*db.Portfolio, error) {
	return s.repo.GetPortfolioByID(userID, portfolioID)
}

func (s *PortfolioService) CreatePortfolio(portfolio *db.Portfolio) (*db.Portfolio, error) {
	return s.repo.CreatePortfolio(portfolio)
}

func (s *PortfolioService) UpdatePortfolio(portfolio *db.Portfolio) (*db.Portfolio, error) {
	return s.repo.UpdatePortfolio(portfolio)
}

func (s *PortfolioService) DeletePortfolio(userID, portfolioID int) error {
	return s.repo.DeletePortfolio(userID, portfolioID)
}

func (s *PortfolioService) GetHoldings(portfolioID int) ([]*db.PortfolioItem, error) {
	return s.repo.GetPortfolioHoldings(portfolioID)
}

func (s *PortfolioService) GetHoldingByID(portfolioID, holdingID int) (*db.PortfolioItem, error) {
	return s.repo.GetHoldingByID(portfolioID, holdingID)
}

func (s *PortfolioService) CreateHolding(item *db.PortfolioItem) (*db.PortfolioItem, error) {
	return s.repo.CreateHolding(item)
}

func (s *PortfolioService) UpdateHolding(item *db.PortfolioItem) (*db.PortfolioItem, error) {
	return s.repo.UpdateHolding(item)
}

func (s *PortfolioService) DeleteHolding(portfolioID, holdingID int) error {
	return s.repo.DeleteHolding(portfolioID, holdingID)
}

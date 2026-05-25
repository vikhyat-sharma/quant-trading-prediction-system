package services

import (
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type WatchlistService struct {
	repo *repositories.WatchlistRepository
}

func NewWatchlistService(repo *repositories.WatchlistRepository) *WatchlistService {
	return &WatchlistService{repo: repo}
}

func (s *WatchlistService) CreateWatchlist(wl *db.UserWatchlist) (*db.UserWatchlist, error) {
	return s.repo.CreateWatchlist(wl)
}

func (s *WatchlistService) GetWatchlistsByUser(userID int) ([]*db.UserWatchlist, error) {
	return s.repo.GetWatchlistsByUser(userID)
}

func (s *WatchlistService) AddStockToWatchlist(item *db.WatchlistItem) (*db.WatchlistItem, error) {
	return s.repo.AddStockToWatchlist(item)
}

func (s *WatchlistService) GetWatchlistItems(watchlistID int) ([]*db.WatchlistItem, error) {
	return s.repo.GetWatchlistItems(watchlistID)
}

func (s *WatchlistService) RemoveStockFromWatchlist(watchlistID, stockID int) error {
	return s.repo.RemoveStockFromWatchlist(watchlistID, stockID)
}

func (s *WatchlistService) DeleteWatchlist(watchlistID int) error {
	return s.repo.DeleteWatchlist(watchlistID)
}

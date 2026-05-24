package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type WatchlistRepository struct {
	db *sql.DB
}

func NewWatchlistRepository(database *sql.DB) *WatchlistRepository {
	return &WatchlistRepository{db: database}
}

func (r *WatchlistRepository) CreateWatchlist(watchlist *db.UserWatchlist) (*db.UserWatchlist, error) {
	var id int
	var createdAt sql.NullTime
	if err := r.db.QueryRow(
		"INSERT INTO user_watchlists (user_id, name) VALUES ($1, $2) RETURNING id, created_at",
		watchlist.UserID, watchlist.Name,
	).Scan(&id, &createdAt); err != nil {
		return nil, err
	}
	watchlist.ID = id
	if createdAt.Valid {
		watchlist.CreatedAt = createdAt.Time
	}
	return watchlist, nil
}

func (r *WatchlistRepository) GetWatchlistsByUser(userID int) ([]*db.UserWatchlist, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, created_at FROM user_watchlists WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lists []*db.UserWatchlist
	for rows.Next() {
		var wl db.UserWatchlist
		if err := rows.Scan(&wl.ID, &wl.UserID, &wl.Name, &wl.CreatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, &wl)
	}
	return lists, nil
}

func (r *WatchlistRepository) AddStockToWatchlist(item *db.WatchlistItem) (*db.WatchlistItem, error) {
	var id int
	var createdAt sql.NullTime
	if err := r.db.QueryRow(
		"INSERT INTO watchlist_items (watchlist_id, stock_id) VALUES ($1, $2) RETURNING id, created_at",
		item.WatchlistID, item.StockID,
	).Scan(&id, &createdAt); err != nil {
		return nil, err
	}
	item.ID = id
	if createdAt.Valid {
		item.CreatedAt = createdAt.Time
	}
	return item, nil
}

func (r *WatchlistRepository) GetWatchlistItems(watchlistID int) ([]*db.WatchlistItem, error) {
	rows, err := r.db.Query("SELECT id, watchlist_id, stock_id, created_at FROM watchlist_items WHERE watchlist_id = $1 ORDER BY created_at DESC", watchlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*db.WatchlistItem
	for rows.Next() {
		var item db.WatchlistItem
		if err := rows.Scan(&item.ID, &item.WatchlistID, &item.StockID, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *WatchlistRepository) RemoveStockFromWatchlist(watchlistID, stockID int) error {
	_, err := r.db.Exec("DELETE FROM watchlist_items WHERE watchlist_id = $1 AND stock_id = $2", watchlistID, stockID)
	return err
}

func (r *WatchlistRepository) DeleteWatchlist(watchlistID int) error {
	_, err := r.db.Exec("DELETE FROM user_watchlists WHERE id = $1", watchlistID)
	return err
}

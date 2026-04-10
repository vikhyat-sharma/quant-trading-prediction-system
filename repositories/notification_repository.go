package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(database *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: database}
}

func (r *NotificationRepository) CreateNotification(notification *db.Notification) (*db.Notification, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO notifications (alert_id, stock_id, price, message) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		notification.AlertID, notification.StockID, notification.Price, notification.Message,
	).Scan(&id, &notification.CreatedAt)
	if err != nil {
		return nil, err
	}

	notification.ID = id
	return notification, nil
}

func (r *NotificationRepository) GetNotificationsByStockID(stockID int) ([]*db.Notification, error) {
	rows, err := r.db.Query(
		"SELECT id, alert_id, stock_id, price, message, created_at FROM notifications WHERE stock_id = $1 ORDER BY created_at DESC",
		stockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*db.Notification
	for rows.Next() {
		var notification db.Notification
		err := rows.Scan(&notification.ID, &notification.AlertID, &notification.StockID, &notification.Price, &notification.Message, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	if len(notifications) == 0 {
		return nil, db.ErrRecordNotFound
	}

	return notifications, nil
}

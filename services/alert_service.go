package services

import (
	"fmt"
	"strings"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type AlertService struct {
	alertRepo        *repositories.AlertRepository
	notificationRepo *repositories.NotificationRepository
	priceHistoryRepo *repositories.PriceHistoryRepository
	stockRepo        *repositories.StockRepository
}

func NewAlertService(
	alertRepo *repositories.AlertRepository,
	notificationRepo *repositories.NotificationRepository,
	priceHistoryRepo *repositories.PriceHistoryRepository,
	stockRepo *repositories.StockRepository,
) *AlertService {
	return &AlertService{
		alertRepo:        alertRepo,
		notificationRepo: notificationRepo,
		priceHistoryRepo: priceHistoryRepo,
		stockRepo:        stockRepo,
	}
}

func (s *AlertService) CreateAlert(alert *db.Alert) (*db.Alert, error) {
	if err := alert.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.stockRepo.GetStock(alert.StockID); err != nil {
		return nil, err
	}

	return s.alertRepo.CreateAlert(alert)
}

func (s *AlertService) GetAlerts(stockID int) ([]*db.Alert, error) {
	return s.alertRepo.GetAlertsByStockID(stockID)
}

func (s *AlertService) DeleteAlert(stockID int, alertID int) error {
	return s.alertRepo.DeleteAlert(stockID, alertID)
}

func (s *AlertService) EvaluateAlerts(stockID int) ([]*db.Notification, error) {
	latest, err := s.priceHistoryRepo.GetLatestPrice(stockID)
	if err != nil {
		return nil, err
	}

	alerts, err := s.alertRepo.GetEnabledAlertsByStockID(stockID)
	if err != nil {
		return nil, err
	}

	var notifications []*db.Notification
	for _, alert := range alerts {
		triggered := false
		message := ""
		condition := strings.ToLower(alert.Condition)

		switch condition {
		case db.AlertConditionAbove:
			if latest.Price >= alert.Threshold {
				triggered = true
				message = fmt.Sprintf("Price alert triggered: %s crossed above %.2f", condition, alert.Threshold)
			}
		case db.AlertConditionBelow:
			if latest.Price <= alert.Threshold {
				triggered = true
				message = fmt.Sprintf("Price alert triggered: %s crossed below %.2f", condition, alert.Threshold)
			}
		}

		if triggered {
			notification := &db.Notification{
				AlertID: alert.ID,
				StockID: stockID,
				Price:   latest.Price,
				Message: message,
			}
			notification, err = s.notificationRepo.CreateNotification(notification)
			if err != nil {
				return nil, err
			}
			notifications = append(notifications, notification)
		}
	}

	return notifications, nil
}

func (s *AlertService) GetNotifications(stockID int) ([]*db.Notification, error) {
	return s.notificationRepo.GetNotificationsByStockID(stockID)
}

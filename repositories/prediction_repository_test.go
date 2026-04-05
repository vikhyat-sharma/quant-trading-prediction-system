package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPredictionRepository_GetPredictionsByStockID_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "stock_id", "predicted_price", "date"}).
		AddRow(1, 1, 150.50, now).
		AddRow(2, 1, 155.75, now)

	mock.ExpectQuery("SELECT id, stock_id, predicted_price, date FROM predictions WHERE stock_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewPredictionRepository(mockDB)
	predictions, err := repo.GetPredictionsByStockID(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(predictions) != 2 {
		t.Errorf("expected 2 predictions, got %d", len(predictions))
	}

	if len(predictions) > 0 && predictions[0].PredictedPrice != 150.50 {
		t.Errorf("expected price 150.50, got %f", predictions[0].PredictedPrice)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPredictionRepository_GetPredictionsByStockID_Empty(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "stock_id", "predicted_price", "date"})

	mock.ExpectQuery("SELECT id, stock_id, predicted_price, date FROM predictions WHERE stock_id = \\$1").
		WithArgs(9999).
		WillReturnRows(rows)

	repo := NewPredictionRepository(mockDB)
	predictions, err := repo.GetPredictionsByStockID(9999)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(predictions) != 0 {
		t.Errorf("expected 0 predictions, got %d", len(predictions))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPredictionRepository_GetPredictionsByStockID_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery("SELECT id, stock_id, predicted_price, date FROM predictions WHERE stock_id = \\$1").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	repo := NewPredictionRepository(mockDB)
	predictions, err := repo.GetPredictionsByStockID(1)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if predictions != nil {
		t.Errorf("expected nil predictions, got %v", predictions)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

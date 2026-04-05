package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStockRepository_GetStock_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "symbol", "name"}).
		AddRow(1, "AAPL", "Apple Inc.")

	mock.ExpectQuery("SELECT id, symbol, name FROM stocks WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewStockRepository(mockDB)
	stock, err := repo.GetStock(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if stock == nil {
		t.Errorf("expected stock, got nil")
	}

	if stock != nil && stock.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", stock.Symbol)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestStockRepository_GetStock_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery("SELECT id, symbol, name FROM stocks WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	repo := NewStockRepository(mockDB)
	stock, err := repo.GetStock(999)

	if err == nil {
		t.Errorf("expected error for non-existent stock, got nil")
	}

	if stock != nil {
		t.Errorf("expected nil stock, got %v", stock)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestStockRepository_GetAllStocks_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "symbol", "name"}).
		AddRow(1, "AAPL", "Apple Inc.").
		AddRow(2, "GOOGL", "Alphabet Inc.")

	mock.ExpectQuery("SELECT id, symbol, name FROM stocks").
		WillReturnRows(rows)

	repo := NewStockRepository(mockDB)
	stocks, err := repo.GetAllStocks()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(stocks) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(stocks))
	}

	if stocks[0].Symbol != "AAPL" {
		t.Errorf("expected first stock symbol AAPL, got %s", stocks[0].Symbol)
	}

	if stocks[1].Symbol != "GOOGL" {
		t.Errorf("expected second stock symbol GOOGL, got %s", stocks[1].Symbol)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestStockRepository_GetAllStocks_Empty(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "symbol", "name"})

	mock.ExpectQuery("SELECT id, symbol, name FROM stocks").
		WillReturnRows(rows)

	repo := NewStockRepository(mockDB)
	stocks, err := repo.GetAllStocks()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(stocks) != 0 {
		t.Errorf("expected 0 stocks, got %d", len(stocks))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

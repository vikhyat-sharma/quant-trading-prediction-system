package algorithms

import (
	"math"
	"testing"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// ── SMA ───────────────────────────────────────────────────────────────────────

func TestCalculateSMA_KnownValues(t *testing.T) {
	prices := []float64{1, 2, 3, 4, 5}
	got := CalculateSMA(prices, 3)
	// SMA(3) of last 3 values [3,4,5] = 4.0
	if !almostEqual(got, 4.0, 1e-9) {
		t.Fatalf("SMA(3) = %f, want 4.0", got)
	}
}

func TestCalculateSMA_InsufficientData(t *testing.T) {
	if got := CalculateSMA([]float64{1, 2}, 5); got != 0 {
		t.Fatalf("expected 0 for insufficient data, got %f", got)
	}
}

func TestCalculateSMA_EmptySlice(t *testing.T) {
	if got := CalculateSMA(nil, 5); got != 0 {
		t.Fatalf("expected 0 for nil slice, got %f", got)
	}
}

// ── EMA ───────────────────────────────────────────────────────────────────────

func TestCalculateEMA_SeedEqualsFirstSMA(t *testing.T) {
	// With exactly `period` prices, EMA should equal SMA.
	prices := []float64{10, 20, 30}
	got := CalculateEMA(prices, 3)
	want := CalculateSMA(prices, 3)
	if !almostEqual(got, want, 1e-9) {
		t.Fatalf("EMA seed = %f, want SMA = %f", got, want)
	}
}

func TestCalculateEMA_InsufficientData(t *testing.T) {
	if got := CalculateEMA([]float64{1, 2}, 5); got != 0 {
		t.Fatalf("expected 0 for insufficient data, got %f", got)
	}
}

// ── RSI ───────────────────────────────────────────────────────────────────────

func TestCalculateRSI_AllGains_Returns100(t *testing.T) {
	prices := make([]float64, 20)
	for i := range prices {
		prices[i] = float64(i + 1)
	}
	got := CalculateRSI(prices, 14)
	if !almostEqual(got, 100.0, 1e-6) {
		t.Fatalf("RSI all-gains = %f, want 100", got)
	}
}

func TestCalculateRSI_AllLosses_Returns0(t *testing.T) {
	prices := make([]float64, 20)
	for i := range prices {
		prices[i] = float64(20 - i)
	}
	got := CalculateRSI(prices, 14)
	if !almostEqual(got, 0.0, 1e-6) {
		t.Fatalf("RSI all-losses = %f, want 0", got)
	}
}

func TestCalculateRSI_InsufficientData_ReturnsNeutral(t *testing.T) {
	got := CalculateRSI([]float64{1, 2, 3}, 14)
	if got != 50 {
		t.Fatalf("RSI insufficient data = %f, want 50", got)
	}
}

func TestCalculateRSI_InRange(t *testing.T) {
	prices := []float64{44.34, 44.09, 44.15, 43.61, 44.33, 44.83, 45.10, 45.15,
		43.61, 44.33, 44.83, 45.10, 45.15, 43.61, 44.33}
	got := CalculateRSI(prices, 14)
	if got < 0 || got > 100 {
		t.Fatalf("RSI out of range [0,100]: %f", got)
	}
}

// ── Bollinger Bands ───────────────────────────────────────────────────────────

func TestCalculateBollingerBands_UpperAboveLower(t *testing.T) {
	prices := make([]float64, 20)
	for i := range prices {
		prices[i] = 100 + float64(i%5)
	}
	upper, lower := CalculateBollingerBands(prices, 20, 2)
	if upper <= lower {
		t.Fatalf("upper band (%f) should be > lower band (%f)", upper, lower)
	}
}

func TestCalculateBollingerBands_InsufficientData(t *testing.T) {
	upper, lower := CalculateBollingerBands([]float64{1, 2}, 20, 2)
	if upper != 0 || lower != 0 {
		t.Fatalf("expected (0,0) for insufficient data, got (%f,%f)", upper, lower)
	}
}

// ── MACD ──────────────────────────────────────────────────────────────────────

func TestCalculateMACD_InsufficientData(t *testing.T) {
	macd, signal := CalculateMACD([]float64{1, 2, 3})
	if macd != 0 || signal != 0 {
		t.Fatalf("expected (0,0) for insufficient data, got (%f,%f)", macd, signal)
	}
}

func TestCalculateMACD_SignalNotHardcodedApproximation(t *testing.T) {
	// With enough data, signal must NOT equal macd*0.7 (the old broken formula).
	prices := make([]float64, 40)
	for i := range prices {
		prices[i] = 100 + float64(i)*0.5
	}
	macd, signal := CalculateMACD(prices)
	if almostEqual(signal, macd*0.7, 1e-6) {
		t.Fatal("signal line appears to use the old hardcoded 0.7 approximation")
	}
}

// ── Prediction strategies ─────────────────────────────────────────────────────

func TestSimpleMovingAveragePrediction_InsufficientData(t *testing.T) {
	result := SimpleMovingAveragePrediction([]float64{1, 2, 3})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.PredictedPrice != 0 {
		t.Fatalf("expected 0 predicted price for insufficient data, got %f", result.PredictedPrice)
	}
	if result.ConfidenceScore != 0.3 {
		t.Fatalf("expected confidence 0.3 for insufficient data, got %f", result.ConfidenceScore)
	}
}

func TestEnsemblePrediction_PositivePriceAndValidConfidence(t *testing.T) {
	prices := make([]float64, 30)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}
	result := EnsemblePrediction(prices)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.PredictedPrice <= 0 {
		t.Fatalf("ensemble predicted price should be > 0, got %f", result.PredictedPrice)
	}
	if result.ConfidenceScore < 0 || result.ConfidenceScore > 1 {
		t.Fatalf("confidence score out of [0,1]: %f", result.ConfidenceScore)
	}
	if result.UpperBound <= result.LowerBound {
		t.Fatalf("upper bound (%f) should exceed lower bound (%f)", result.UpperBound, result.LowerBound)
	}
}

func TestEnsemblePrediction_InsufficientData(t *testing.T) {
	result := EnsemblePrediction([]float64{1, 2, 3})
	if result == nil || result.PredictedPrice != 0 {
		t.Fatalf("expected zero-price result for insufficient data")
	}
}

// ── Backtest ──────────────────────────────────────────────────────────────────

func TestBacktestStrategy_NilAlgorithm(t *testing.T) {
	result := BacktestStrategy([]float64{1, 2, 3}, nil)
	if result == nil {
		t.Fatal("expected non-nil result for nil algorithm")
	}
	if result.Trades != 0 {
		t.Fatalf("expected 0 trades for nil algorithm, got %d", result.Trades)
	}
}

func TestBacktestStrategy_InsufficientPrices(t *testing.T) {
	result := BacktestStrategy([]float64{100}, SimpleMovingAveragePrediction)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestBacktestStrategy_WinRateInRange(t *testing.T) {
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100 + float64(i%5)
	}
	result := BacktestStrategy(prices, EnsemblePrediction)
	if result.WinRate < 0 || result.WinRate > 100 {
		t.Fatalf("win rate out of [0,100]: %f", result.WinRate)
	}
}

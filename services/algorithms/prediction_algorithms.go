package algorithms

import (
	"math"
)

// PredictionResult holds the result of a prediction algorithm
type PredictionResult struct {
	PredictedPrice  float64
	ConfidenceScore float64
	UpperBound      float64
	LowerBound      float64
	Algorithm       string
}

type BacktestResult struct {
	Strategy         string    `json:"strategy"`
	TotalReturn      float64   `json:"total_return"`
	AnnualizedReturn float64   `json:"annualized_return"`
	WinRate          float64   `json:"win_rate"`
	MaxDrawdown      float64   `json:"max_drawdown"`
	Trades           int       `json:"trades"`
	PositiveTrades   int       `json:"positive_trades"`
	NegativeTrades   int       `json:"negative_trades"`
	EquityCurve      []float64 `json:"equity_curve"`
	DailyReturns     []float64 `json:"daily_returns"`
}

// TechnicalIndicators holds calculated technical indicators
type TechnicalIndicators struct {
	SMA20      float64 // 20-day Simple Moving Average
	EMA12      float64 // 12-day Exponential Moving Average
	RSI        float64 // Relative Strength Index
	MACD       float64 // MACD line
	MACDSignal float64 // MACD signal line
	BBUpper    float64 // Bollinger Band Upper
	BBLower    float64 // Bollinger Band Lower
}

// CalculateSMA calculates Simple Moving Average
func CalculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// CalculateEMA calculates Exponential Moving Average
func CalculateEMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	multiplier := 2.0 / (float64(period) + 1)

	// Start with SMA
	ema := CalculateSMA(prices[:period], period)

	// Calculate EMA for remaining prices
	for i := period; i < len(prices); i++ {
		ema = (prices[i] * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

// CalculateRSI calculates the Relative Strength Index using Wilder's smoothed
// moving average method (the standard definition).
// Requires at least period+1 price points.
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // Neutral RSI when insufficient data
	}

	// Seed: simple average of first `period` changes
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain += change
		} else {
			avgLoss += -change
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Wilder's smoothing for remaining prices
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + (-change)) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// CalculateMACD calculates MACD (Moving Average Convergence Divergence).
// MACD line = EMA(12) - EMA(26)
// Signal line = 9-period EMA of the MACD line.
// Requires at least 26 + 9 - 1 = 34 price points for a meaningful signal.
func CalculateMACD(prices []float64) (macd, signal float64) {
	if len(prices) < 26 {
		return 0, 0
	}
	ema12 := CalculateEMA(prices, 12)
	ema26 := CalculateEMA(prices, 26)
	macd = ema12 - ema26

	// Build a MACD series over the available history to compute the signal EMA.
	// We need at least 26 points to get the first MACD value, then 9 more for signal.
	if len(prices) < 34 {
		// Not enough data for a proper signal line; return MACD only.
		signal = macd
		return
	}

	macdSeries := make([]float64, 0, len(prices)-25)
	for i := 26; i <= len(prices); i++ {
		e12 := CalculateEMA(prices[:i], 12)
		e26 := CalculateEMA(prices[:i], 26)
		macdSeries = append(macdSeries, e12-e26)
	}
	signal = CalculateEMA(macdSeries, 9)
	return
}

// CalculateBollingerBands calculates Bollinger Bands
func CalculateBollingerBands(prices []float64, period int, stdDev float64) (upper, lower float64) {
	if len(prices) < period {
		return 0, 0
	}

	sma := CalculateSMA(prices, period)

	// Calculate standard deviation
	variance := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		variance += math.Pow(prices[i]-sma, 2)
	}
	variance /= float64(period)
	std := math.Sqrt(variance)

	upper = sma + (std * stdDev)
	lower = sma - (std * stdDev)

	return upper, lower
}

// SimpleMovingAveragePrediction predicts price using SMA
func SimpleMovingAveragePrediction(prices []float64) *PredictionResult {
	if len(prices) < 20 {
		return &PredictionResult{
			Algorithm:       "SMA",
			ConfidenceScore: 0.3,
		}
	}

	sma20 := CalculateSMA(prices, 20)
	sma10 := CalculateSMA(prices, 10)
	currentPrice := prices[len(prices)-1]

	// Trend determination
	trend := sma10 - sma20
	predictedPrice := currentPrice + trend

	confidence := 0.6
	if math.Abs(trend/currentPrice) < 0.02 {
		confidence = 0.5 // Low confidence if trend is small
	}

	variation := currentPrice * 0.05 // 5% variation band
	return &PredictionResult{
		PredictedPrice:  predictedPrice,
		ConfidenceScore: confidence,
		UpperBound:      predictedPrice + variation,
		LowerBound:      predictedPrice - variation,
		Algorithm:       "SMA",
	}
}

// ExponentialMovingAveragePrediction predicts price using EMA
func ExponentialMovingAveragePrediction(prices []float64) *PredictionResult {
	if len(prices) < 12 {
		return &PredictionResult{
			Algorithm:       "EMA",
			ConfidenceScore: 0.3,
		}
	}

	ema12 := CalculateEMA(prices, 12)
	ema5 := CalculateEMA(prices, 5)
	currentPrice := prices[len(prices)-1]

	// Trend determination
	trend := ema5 - ema12
	predictedPrice := currentPrice + (trend * 0.8)

	confidence := 0.65
	if math.Abs(trend/currentPrice) < 0.01 {
		confidence = 0.55
	}

	variation := currentPrice * 0.06 // 6% variation band
	return &PredictionResult{
		PredictedPrice:  predictedPrice,
		ConfidenceScore: confidence,
		UpperBound:      predictedPrice + variation,
		LowerBound:      predictedPrice - variation,
		Algorithm:       "EMA",
	}
}

// MomentumPrediction predicts price using momentum and RSI
func MomentumPrediction(prices []float64) *PredictionResult {
	if len(prices) < 14 {
		return &PredictionResult{
			Algorithm:       "MOMENTUM",
			ConfidenceScore: 0.3,
		}
	}

	rsi := CalculateRSI(prices, 14)
	currentPrice := prices[len(prices)-1]

	// Momentum calculation
	momentum := prices[len(prices)-1] - prices[len(prices)-5]
	predictedPrice := currentPrice + momentum

	// Adjust based on RSI
	confidence := 0.55
	if rsi > 70 {
		confidence = 0.6                                 // Overbought signal
		predictedPrice = currentPrice - (momentum * 0.5) // Reversal expected
	} else if rsi < 30 {
		confidence = 0.6                                 // Oversold signal
		predictedPrice = currentPrice + (momentum * 1.2) // Recovery expected
	}

	variation := currentPrice * 0.07 // 7% variation band
	return &PredictionResult{
		PredictedPrice:  predictedPrice,
		ConfidenceScore: confidence,
		UpperBound:      predictedPrice + variation,
		LowerBound:      predictedPrice - variation,
		Algorithm:       "MOMENTUM",
	}
}

// MeanReversionPrediction predicts price based on mean reversion
func MeanReversionPrediction(prices []float64) *PredictionResult {
	if len(prices) < 30 {
		return &PredictionResult{
			Algorithm:       "MEAN_REVERSION",
			ConfidenceScore: 0.3,
		}
	}

	mean := 0.0
	for _, p := range prices {
		mean += p
	}
	mean /= float64(len(prices))

	currentPrice := prices[len(prices)-1]
	deviation := currentPrice - mean
	deviationPercent := (deviation / mean) * 100

	// If price deviates significantly, expect reversion
	confidence := 0.5
	predictedPrice := currentPrice

	if math.Abs(deviationPercent) > 5 {
		// Strong deviation - expect reversion
		confidence = 0.65
		predictedPrice = mean + (deviation * 0.3) // Revert partially
	} else if math.Abs(deviationPercent) > 2 {
		confidence = 0.58
		predictedPrice = mean + (deviation * 0.5)
	}

	variation := currentPrice * 0.08 // 8% variation band
	return &PredictionResult{
		PredictedPrice:  predictedPrice,
		ConfidenceScore: confidence,
		UpperBound:      predictedPrice + variation,
		LowerBound:      predictedPrice - variation,
		Algorithm:       "MEAN_REVERSION",
	}
}

// EnsemblePrediction combines multiple algorithms for a weighted prediction.
// Sub-strategies that return a zero predicted price (insufficient data) are
// excluded from the weighted average to avoid pulling the result toward zero.
func EnsemblePrediction(prices []float64) *PredictionResult {
	if len(prices) < 30 {
		return &PredictionResult{
			Algorithm:       "ENSEMBLE",
			ConfidenceScore: 0.3,
		}
	}

	type weightedStrategy struct {
		result *PredictionResult
		weight float64
	}

	candidates := []weightedStrategy{
		{SimpleMovingAveragePrediction(prices), 0.25},
		{ExponentialMovingAveragePrediction(prices), 0.30},
		{MomentumPrediction(prices), 0.25},
		{MeanReversionPrediction(prices), 0.20},
	}

	weightedPrice := 0.0
	weightedConfidence := 0.0
	totalWeight := 0.0
	weightedVariation := 0.0

	for _, c := range candidates {
		if c.result == nil || c.result.PredictedPrice <= 0 {
			continue
		}
		weightedPrice += c.result.PredictedPrice * c.weight
		weightedConfidence += c.result.ConfidenceScore * c.weight
		weightedVariation += ((c.result.UpperBound - c.result.LowerBound) / 2) * c.weight
		totalWeight += c.weight
	}

	if totalWeight == 0 {
		return &PredictionResult{Algorithm: "ENSEMBLE", ConfidenceScore: 0.3}
	}

	predictedPrice := weightedPrice / totalWeight
	confidence := weightedConfidence / totalWeight
	variation := weightedVariation / totalWeight

	return &PredictionResult{
		PredictedPrice:  predictedPrice,
		ConfidenceScore: confidence,
		UpperBound:      predictedPrice + variation,
		LowerBound:      predictedPrice - variation,
		Algorithm:       "ENSEMBLE",
	}
}

func BacktestStrategy(prices []float64, algorithm func([]float64) *PredictionResult) *BacktestResult {
	result := &BacktestResult{
		Strategy:     "UNKNOWN",
		EquityCurve:  make([]float64, 0),
		DailyReturns: make([]float64, 0),
	}

	if algorithm == nil || len(prices) < 2 {
		return result
	}

	firstPrediction := algorithm(prices[:1])
	if firstPrediction != nil {
		result.Strategy = firstPrediction.Algorithm
	}

	equity := 1.0
	peak := 1.0
	positiveTrades := 0
	negativeTrades := 0
	trades := 0

	for i := 1; i < len(prices); i++ {
		window := prices[:i]
		currentPrice := prices[i-1]
		nextPrice := prices[i]

		prediction := algorithm(window)
		if prediction == nil {
			result.EquityCurve = append(result.EquityCurve, equity)
			result.DailyReturns = append(result.DailyReturns, 0)
			continue
		}

		tradeReturn := 0.0
		if prediction.PredictedPrice > currentPrice {
			tradeReturn = (nextPrice - currentPrice) / currentPrice
			trades++
			if tradeReturn > 0 {
				positiveTrades++
			} else if tradeReturn < 0 {
				negativeTrades++
			}
		}

		equity *= 1 + tradeReturn
		if equity > peak {
			peak = equity
		}
		maxDrawdown := (peak - equity) / peak
		if maxDrawdown > result.MaxDrawdown {
			result.MaxDrawdown = maxDrawdown
		}

		result.EquityCurve = append(result.EquityCurve, equity)
		result.DailyReturns = append(result.DailyReturns, tradeReturn)
	}

	result.Trades = trades
	result.PositiveTrades = positiveTrades
	result.NegativeTrades = negativeTrades
	if trades > 0 {
		result.WinRate = float64(positiveTrades) / float64(trades) * 100
	}
	result.TotalReturn = equity - 1
	if len(result.DailyReturns) > 0 {
		periods := float64(len(result.DailyReturns))
		result.AnnualizedReturn = math.Pow(1+result.TotalReturn, 252.0/periods) - 1
	}

	return result
}

// CalculateTechnicalIndicators calculates all technical indicators
func CalculateTechnicalIndicators(prices []float64) *TechnicalIndicators {
	if len(prices) < 26 {
		return &TechnicalIndicators{}
	}

	sma20 := CalculateSMA(prices, 20)
	ema12 := CalculateEMA(prices, 12)
	rsi := CalculateRSI(prices, 14)
	macd, signal := CalculateMACD(prices)
	bbUpper, bbLower := CalculateBollingerBands(prices, 20, 2)

	return &TechnicalIndicators{
		SMA20:      sma20,
		EMA12:      ema12,
		RSI:        rsi,
		MACD:       macd,
		MACDSignal: signal,
		BBUpper:    bbUpper,
		BBLower:    bbLower,
	}
}

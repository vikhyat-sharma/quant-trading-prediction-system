package services

import (
	"math"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type PortfolioAnalyticsService struct {
	portfolioRepo    *repositories.PortfolioRepository
	priceHistoryRepo *repositories.PriceHistoryRepository
}

func NewPortfolioAnalyticsService(
	portfolioRepo *repositories.PortfolioRepository,
	priceHistoryRepo *repositories.PriceHistoryRepository,
) *PortfolioAnalyticsService {
	return &PortfolioAnalyticsService{
		portfolioRepo:    portfolioRepo,
		priceHistoryRepo: priceHistoryRepo,
	}
}

// CalculatePortfolioValue calculates the current value of a portfolio
func (s *PortfolioAnalyticsService) CalculatePortfolioValue(portfolioID int) (float64, error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return 0, err
	}

	totalValue := 0.0
	for _, holding := range holdings {
		latestPrice, err := s.priceHistoryRepo.GetLatestPrice(holding.StockID)
		if err == nil && latestPrice != nil {
			totalValue += latestPrice.Price * holding.Quantity
		}
	}

	return totalValue, nil
}

// CalculatePortfolioMetrics calculates gains/losses and returns for a portfolio
func (s *PortfolioAnalyticsService) CalculatePortfolioMetrics(portfolioID int) (totalValue, costBasis, gainLoss, returnPercent float64, err error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return
	}

	totalValue = 0.0
	costBasis = 0.0

	for _, holding := range holdings {
		// Add to cost basis
		costBasis += holding.AvgCost * holding.Quantity

		// Get current price and add to total value
		latestPrice, err := s.priceHistoryRepo.GetLatestPrice(holding.StockID)
		if err == nil && latestPrice != nil {
			totalValue += latestPrice.Price * holding.Quantity
		} else {
			// If no price data, use cost basis
			totalValue += holding.AvgCost * holding.Quantity
		}
	}

	gainLoss = totalValue - costBasis
	if costBasis > 0 {
		returnPercent = (gainLoss / costBasis) * 100
	}

	return
}

// CalculatePortfolioVolatility calculates the volatility of portfolio returns
func (s *PortfolioAnalyticsService) CalculatePortfolioVolatility(portfolioID int, days int) (float64, error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return 0, err
	}

	if len(holdings) == 0 {
		return 0, nil
	}

	returns := make([]float64, 0)

	for i := 0; i < days; i++ {
		endDate := time.Now().AddDate(0, 0, -i)
		startDate := endDate.AddDate(0, 0, -1)

		dayValue := 0.0
		prevDayValue := 0.0

		for _, holding := range holdings {
			// Get price for current day
			priceHist, _ := s.priceHistoryRepo.GetPriceHistoryByStockIDAndDateRange(
				holding.StockID, startDate, endDate,
			)
			if len(priceHist) > 0 {
				currentPrice := priceHist[len(priceHist)-1].Price
				dayValue += currentPrice * holding.Quantity
				if len(priceHist) > 1 {
					prevPrice := priceHist[0].Price
					prevDayValue += prevPrice * holding.Quantity
				}
			}
		}

		if prevDayValue > 0 && dayValue > 0 {
			dailyReturn := (dayValue - prevDayValue) / prevDayValue
			returns = append(returns, dailyReturn)
		}
	}

	if len(returns) == 0 {
		return 0, nil
	}

	// Calculate standard deviation of returns
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns))

	volatility := math.Sqrt(variance) * math.Sqrt(252) // Annualized volatility

	return volatility, nil
}

// CalculateSharpeRatio calculates Sharpe ratio for the portfolio
func (s *PortfolioAnalyticsService) CalculateSharpeRatio(portfolioID int, days int, riskFreeRate float64) (float64, error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return 0, err
	}

	if len(holdings) == 0 {
		return 0, nil
	}

	returns := make([]float64, 0)
	values := make([]float64, 0)

	for i := 0; i < days; i++ {
		endDate := time.Now().AddDate(0, 0, -i)
		startDate := endDate.AddDate(0, 0, -1)

		dayValue := 0.0
		prevDayValue := 0.0

		for _, holding := range holdings {
			priceHist, _ := s.priceHistoryRepo.GetPriceHistoryByStockIDAndDateRange(
				holding.StockID, startDate, endDate,
			)
			if len(priceHist) > 0 {
				currentPrice := priceHist[len(priceHist)-1].Price
				dayValue += currentPrice * holding.Quantity
				if len(priceHist) > 1 {
					prevPrice := priceHist[0].Price
					prevDayValue += prevPrice * holding.Quantity
				}
			}
		}

		if prevDayValue > 0 && dayValue > 0 {
			dailyReturn := (dayValue - prevDayValue) / prevDayValue
			returns = append(returns, dailyReturn)
			values = append(values, dayValue)
		}
	}

	if len(returns) == 0 {
		return 0, nil
	}

	// Calculate mean return
	meanReturn := 0.0
	for _, r := range returns {
		meanReturn += r
	}
	meanReturn /= float64(len(returns))
	annualizedReturn := meanReturn * 252

	// Calculate standard deviation
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-meanReturn, 2)
	}
	variance /= float64(len(returns))
	stdDev := math.Sqrt(variance) * math.Sqrt(252)

	// Calculate Sharpe ratio
	sharpe := (annualizedReturn - riskFreeRate) / stdDev

	return sharpe, nil
}

// GetPortfolioPerformanceHistory gets historical performance records for a portfolio
func (s *PortfolioAnalyticsService) GetPortfolioPerformanceHistory(portfolioID int, days int) ([]*db.PortfolioPerformance, error) {
	performances := make([]*db.PortfolioPerformance, 0)

	for i := 0; i < days; i++ {
		recordDate := time.Now().AddDate(0, 0, -i)

		totalValue, costBasis, gainLoss, returnPercent, err := s.CalculatePortfolioMetrics(portfolioID)
		if err != nil {
			continue
		}

		volatility, _ := s.CalculatePortfolioVolatility(portfolioID, 30)
		sharpe, _ := s.CalculateSharpeRatio(portfolioID, 30, 0.02)

		performance := &db.PortfolioPerformance{
			PortfolioID:   portfolioID,
			TotalValue:    totalValue,
			CostBasis:     costBasis,
			GainLoss:      gainLoss,
			ReturnPercent: returnPercent,
			Volatility:    volatility,
			Sharpe:        sharpe,
			RecordDate:    recordDate,
		}

		performances = append(performances, performance)
	}

	return performances, nil
}

// CalculateDiversificationScore calculates how diversified the portfolio is (0-1)
func (s *PortfolioAnalyticsService) CalculateDiversificationScore(portfolioID int) (float64, error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return 0, err
	}

	if len(holdings) == 0 {
		return 0, nil
	}

	// Calculate portfolio value
	totalValue := 0.0
	for _, holding := range holdings {
		latestPrice, err := s.priceHistoryRepo.GetLatestPrice(holding.StockID)
		if err == nil && latestPrice != nil {
			totalValue += latestPrice.Price * holding.Quantity
		} else {
			totalValue += holding.AvgCost * holding.Quantity
		}
	}

	if totalValue == 0 {
		return 0, nil
	}

	// Calculate weights and Herfindahl index
	herfindahl := 0.0
	for _, holding := range holdings {
		holdingValue := 0.0
		latestPrice, err := s.priceHistoryRepo.GetLatestPrice(holding.StockID)
		if err == nil && latestPrice != nil {
			holdingValue = latestPrice.Price * holding.Quantity
		} else {
			holdingValue = holding.AvgCost * holding.Quantity
		}

		weight := holdingValue / totalValue
		herfindahl += weight * weight
	}

	// Convert to diversification score (1 is fully diversified, 0 is concentrated)
	diversificationScore := (1.0 - herfindahl) / (1.0 - (1.0 / float64(len(holdings))))

	// Clamp to 0-1 range
	if diversificationScore < 0 {
		diversificationScore = 0
	}
	if diversificationScore > 1 {
		diversificationScore = 1
	}

	return diversificationScore, nil
}

// GetTopHoldings returns the top N holdings by value
func (s *PortfolioAnalyticsService) GetTopHoldings(portfolioID int, topN int) ([]map[string]interface{}, error) {
	holdings, err := s.portfolioRepo.GetPortfolioHoldings(portfolioID)
	if err != nil {
		return nil, err
	}

	type holding struct {
		Item  *db.PortfolioItem
		Value float64
	}

	holdings_with_values := make([]holding, 0)

	for _, h := range holdings {
		latestPrice, err := s.priceHistoryRepo.GetLatestPrice(h.StockID)
		if err == nil && latestPrice != nil {
			holdings_with_values = append(holdings_with_values, holding{
				Item:  h,
				Value: latestPrice.Price * h.Quantity,
			})
		}
	}

	// Sort by value (descending)
	for i := 0; i < len(holdings_with_values)-1; i++ {
		for j := i + 1; j < len(holdings_with_values); j++ {
			if holdings_with_values[j].Value > holdings_with_values[i].Value {
				holdings_with_values[i], holdings_with_values[j] = holdings_with_values[j], holdings_with_values[i]
			}
		}
	}

	// Get top N
	result := make([]map[string]interface{}, 0)
	limit := topN
	if limit > len(holdings_with_values) {
		limit = len(holdings_with_values)
	}

	for i := 0; i < limit; i++ {
		result = append(result, map[string]interface{}{
			"stock_id": holdings_with_values[i].Item.StockID,
			"quantity": holdings_with_values[i].Item.Quantity,
			"avg_cost": holdings_with_values[i].Item.AvgCost,
			"value":    holdings_with_values[i].Value,
			"weight":   0, // Will be calculated by caller if needed
		})
	}

	return result, nil
}

package services

import "testing"

func TestSentimentService_Analyze_PositiveText(t *testing.T) {
	service := NewSentimentService()
	result := service.Analyze("The market is bullish and the stock is showing strong momentum with a breakout.")

	if result == nil {
		t.Fatal("expected sentiment result, got nil")
	}

	if result.Score <= 0 {
		t.Errorf("expected positive sentiment score, got %f", result.Score)
	}

	if result.Label != "positive" {
		t.Errorf("expected label positive, got %s", result.Label)
	}
}

func TestSentimentService_Analyze_NegativeText(t *testing.T) {
	service := NewSentimentService()
	result := service.Analyze("The stock is facing risk and loss with a drop in momentum and bearish sentiment.")

	if result == nil {
		t.Fatal("expected sentiment result, got nil")
	}

	if result.Score >= 0 {
		t.Errorf("expected negative sentiment score, got %f", result.Score)
	}

	if result.Label != "negative" {
		t.Errorf("expected label negative, got %s", result.Label)
	}
}

func TestSentimentService_Analyze_NeutralText(t *testing.T) {
	service := NewSentimentService()
	result := service.Analyze("The company released earnings and the market reacted with mixed headlines.")

	if result == nil {
		t.Fatal("expected sentiment result, got nil")
	}

	if result.Label != "neutral" {
		t.Errorf("expected label neutral, got %s", result.Label)
	}
}

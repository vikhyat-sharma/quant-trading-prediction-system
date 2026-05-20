package services

import (
	"math"
	"regexp"
	"strings"
)

type SentimentService struct {
	positiveWords map[string]struct{}
	negativeWords map[string]struct{}
}

type SentimentAnalysisResult struct {
	Text          string   `json:"text"`
	Score         float64  `json:"score"`
	Label         string   `json:"label"`
	PositiveCount int      `json:"positive_count"`
	NegativeCount int      `json:"negative_count"`
	PositiveTerms []string `json:"positive_terms"`
	NegativeTerms []string `json:"negative_terms"`
}

func NewSentimentService() *SentimentService {
	return &SentimentService{
		positiveWords: map[string]struct{}{
			"bull": {}, "bullish": {}, "gain": {}, "gains": {}, "up": {}, "positive": {}, "strong": {}, "surge": {}, "beat": {}, "beating": {}, "record": {}, "breakout": {}, "optimism": {}, "optimistic": {}, "rebound": {}, "support": {}, "buy": {}, "buys": {}, "rally": {}, "rallied": {}, "momentum": {}, "growth": {}, "profit": {}, "profits": {}, "upgrade": {}, "favorable": {}, "bullishness": {},
		},
		negativeWords: map[string]struct{}{
			"bear": {}, "bearish": {}, "loss": {}, "losses": {}, "down": {}, "negative": {}, "weak": {}, "drop": {}, "dropped": {}, "miss": {}, "missed": {}, "pullback": {}, "decline": {}, "fall": {}, "falling": {}, "sell": {}, "sells": {}, "recession": {}, "risk": {}, "selloff": {}, "weakness": {}, "cut": {}, "cuts": {}, "downgrade": {}, "uncertain": {}, "uncertainty": {}, "fear": {},
		},
	}
}

func (s *SentimentService) Analyze(text string) *SentimentAnalysisResult {
	result := &SentimentAnalysisResult{
		Text:          strings.TrimSpace(text),
		PositiveTerms: []string{},
		NegativeTerms: []string{},
	}

	cleanText := strings.ToLower(result.Text)
	cleanText = regexp.MustCompile(`[^a-z0-9\s]`).ReplaceAllString(cleanText, " ")

tokens := strings.Fields(cleanText)
	for _, token := range tokens {
		if token == "" {
			continue
		}

		if _, ok := s.positiveWords[token]; ok {
			result.PositiveCount++
			result.PositiveTerms = append(result.PositiveTerms, token)
			continue
		}

		if _, ok := s.negativeWords[token]; ok {
			result.NegativeCount++
			result.NegativeTerms = append(result.NegativeTerms, token)
		}
	}

	totalTerms := len(tokens)
	if totalTerms == 0 {
		result.Score = 0
		result.Label = "neutral"
		return result
	}

	rawScore := float64(result.PositiveCount-result.NegativeCount) / math.Max(float64(totalTerms), 1)
	result.Score = math.Max(-1.0, math.Min(1.0, rawScore))

	switch {
	case result.Score >= 0.15:
		result.Label = "positive"
	case result.Score <= -0.15:
		result.Label = "negative"
	default:
		result.Label = "neutral"
	}

	return result
}

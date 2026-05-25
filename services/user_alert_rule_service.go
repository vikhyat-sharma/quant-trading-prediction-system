package services

import (
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type UserAlertRuleService struct {
	repo *repositories.UserAlertRuleRepository
}

func NewUserAlertRuleService(repo *repositories.UserAlertRuleRepository) *UserAlertRuleService {
	return &UserAlertRuleService{repo: repo}
}

func (s *UserAlertRuleService) CreateAlertRule(rule *db.UserAlertRule) (*db.UserAlertRule, error) {
	return s.repo.CreateAlertRule(rule)
}

func (s *UserAlertRuleService) GetAlertRulesByUser(userID int) ([]*db.UserAlertRule, error) {
	return s.repo.GetAlertRulesByUser(userID)
}

func (s *UserAlertRuleService) DeleteAlertRule(ruleID int) error {
	return s.repo.DeleteAlertRule(ruleID)
}

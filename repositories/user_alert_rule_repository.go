package repositories

import (
	"database/sql"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type UserAlertRuleRepository struct {
	db *sql.DB
}

func NewUserAlertRuleRepository(database *sql.DB) *UserAlertRuleRepository {
	return &UserAlertRuleRepository{db: database}
}

func (r *UserAlertRuleRepository) CreateAlertRule(rule *db.UserAlertRule) (*db.UserAlertRule, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO user_alert_rules (user_id, stock_id, threshold, condition, enabled) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		rule.UserID, rule.StockID, rule.Threshold, rule.Condition, rule.Enabled,
	).Scan(&id, &rule.CreatedAt); err != nil {
		return nil, err
	}
	rule.ID = id
	return rule, nil
}

func (r *UserAlertRuleRepository) GetAlertRulesByUser(userID int) ([]*db.UserAlertRule, error) {
	rows, err := r.db.Query("SELECT id, user_id, stock_id, threshold, condition, enabled, created_at FROM user_alert_rules WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []*db.UserAlertRule
	for rows.Next() {
		var rule db.UserAlertRule
		if err := rows.Scan(&rule.ID, &rule.UserID, &rule.StockID, &rule.Threshold, &rule.Condition, &rule.Enabled, &rule.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (r *UserAlertRuleRepository) DeleteAlertRule(ruleID int) error {
	_, err := r.db.Exec("DELETE FROM user_alert_rules WHERE id = $1", ruleID)
	return err
}

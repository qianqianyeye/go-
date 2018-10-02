package model
type RelActivityRuleinfo struct {
	ID int64 `gorm:"column:id" json:"id"`
	RuleName string `gorm:"column:rule_name" json:"rule_name"`
	RuleType int64 `gorm:"column:rule_type" json:"rule_type"`
	RuleStatus int64 `gorm:"column:rule_status" json:"rule_status"`
}

func (RelActivityRuleinfo) TableName() string {
	return "rel_activity_ruleinfo"
}
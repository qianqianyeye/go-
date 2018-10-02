package model
type RelActivityRule struct {
	ID int64 `gorm:"column:id" json:"id"`
	MerchantId int64 `gorm:"column:_merchant_id" json:"merchant_id"`
	StoreId int64 `gorm:"column:_store_id" json:"store_id"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	Body string `gorm:"column:body" json:"body"`
	RuleId int64 `gorm:"column:rule_id" json:"rule_id"`
	RelActivityRuleinfo RelActivityRuleinfo `gorm:"ForeignKey:RuleId;AssociationForeignKey:;ID" json:"rel_activity_ruleinfo,omitempty"`
}
func (RelActivityRule) TableName() string {
	return "rel_activity_rule"
}

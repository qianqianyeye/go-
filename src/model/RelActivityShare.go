package model
type RelActivityShare struct {
	ID int64 `gorm:"column:id" json:"id"`
	GiveMinMcoin int64 `gorm:"column:give_min_mcoin" json:"give_min_mcoin"`
	GiveMaxMcoin int64 `gorm:"column:give_max_mcoin" json:"give_max_mcoin"`
	MaxShareCoin int64 `gorm:"column:max_share_coin" json:"max_share_coin"`
	DayMaxShareCoin int64 `gorm:"column:day_max_share_coin" json:"day_max_share_coin"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
}

func (RelActivityShare) TableName() string {
	return "rel_activity_share"
}
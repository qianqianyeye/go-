package model
type RelActivityThresholdinfo struct {
	ID int64 `gorm:"column:id" json:"id"`
	ThresholdName string `gorm:"column:threshold_name" json:"threshold_name"`
	ThresholdType int64 `gorm:"column:threshold_type" json:"threshold_type"`
	ThresholdStatus int64 `gorm:"column:threshold_status" json:"threshold_status"`
}

func (RelActivityThresholdinfo) TableName() string {
	return "rel_activity_thresholdinfo"
}
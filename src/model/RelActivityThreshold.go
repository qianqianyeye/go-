package model
type RelActivityThreshold struct {
	ID int64 `gorm:"column:id" json:"id"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	ThresholdId int64 `gorm:"column:threshold_id" json:"threshold_id"`
	Body string `gorm:"column:body" json:"body"`
	RelActivityThresholdinfo  RelActivityThresholdinfo`gorm:"ForeignKey:ThresholdId;AssociationForeignKey:;ID" json:"rel_activity_thresholdinfo,omitempty"`
}

func (RelActivityThreshold) TableName() string {
	return "rel_activity_threshold"
}
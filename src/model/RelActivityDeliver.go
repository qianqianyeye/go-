package model

type RelActivityDeliver struct {
	ID int64 `gorm:"column:id" json:"id"`
	DeliverId int64 `gorm:"column:deliver_id" json:"deliver_id"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	Deliver Deliver `gorm:"ForeignKey:ID;AssociationForeignKey:DeliverId;" json:"Deliver"`
}
func (RelActivityDeliver) TableName() string {
	return "rel_activity_deliver"
}
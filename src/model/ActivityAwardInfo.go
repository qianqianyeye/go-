package model

type ActivityAwardInfo struct {
	ID int64 `gorm:"column:id" json:"id"`
	GoodsId int64 `gorm:"column:goods_id" json:"goods_id" binding:"required"`
	Lv int64 `gorm:"column:lv" json:"lv"`
	Number int64 `gorm:"column:number" json:"number"  binding:"required"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	Goods Goods `gorm:"ForeignKey:ID;AssociationForeignKey:GoodsId;" json:"goods,omitempty"`
	//RelDelverGoods []RelDelverGoods `gorm:"ForeignKey:GoodsId;AssociationForeignKey:GoodsId;" json:"-"`
	RelActivityDeliver []RelActivityDeliver `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ActivityId;" json:"-"`
}
func (ActivityAwardInfo) TableName() string {
	return "activity_award_info"
}
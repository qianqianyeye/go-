package model
type RelDelverGoods struct {
	ID int64 `gorm:"column:id" json:"id"`
	DeliverId int64 `gorm:"column:deliver_id" json:"deliver_id"`
	GoodsId int64 `gorm:"column:goods_id" json:"goods_id"`
	Nums int64 `gorm:"column:nums" json:"nums"`
	Deliver Deliver `gorm:"ForeignKey:ID;AssociationForeignKey:DeliverId;" json:"-"`
}
func (RelDelverGoods) TableName() string {
	return "rel_delver_goods"
}
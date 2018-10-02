package model

import (
	"time"
)

type RelActivityOpenluck struct {
	ID int64 `gorm:"column:id" json:"id"`
	MemberId int64 `gorm:"column:member_id" json:"member_id"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	MerchantId int64 `gorm:"column:_merchant_id" json:"merchant_id"`
	StoreId int64 `gorm:"column:_store_id" json:"store_id"`
	OpenTime time.Time `gorm:"column:open_time" json:"open_time"`
	Activity Activity `gorm:"ForeignKey:id;AssociationForeignKey:ActivityId;" json:"activity,omitempty"`
	LogActivityUserluck []LogActivityUserluck `gorm:"ForeignKey:memberId;AssociationForeignKey:memberId;" json:"log_activity_userluck,omitempty"`
	Member Member `gorm:"ForeignKey:memberId;AssociationForeignKey:id;" json:"member,omitempty"`

	Deliver map[string]interface{} `gorm:"-" json:"deliver,omitempty"`
}
func (RelActivityOpenluck) TableName() string {
	return "rel_activity_openluck"
}
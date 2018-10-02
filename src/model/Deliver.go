package model

import "time"

type Deliver struct {
	ID int64 `gorm:"column:id" json:"id"`
	//Type int64 `gorm:"column:type" json:"type"`
	Address string `gorm:"column:address" json:"address"`
	Status int64 `gorm:"column:status" json:"status"`
	Consignee string `gorm:"column:consignee" json:"consignee"`
	Mobile int64 `gorm:"column:mobile" json:"mobile"`
	MerchantId int64 `gorm:"column:merchant_id" json:"merchant_id"`
	MemberId int64 `gorm:"column:member_id" json:"member_id"`
	ExpressName string `gorm:"column:express_name" json:"express_name"`
	ExpressNo string `gorm:"column:express_no" json:"express_no"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (Deliver) TableName() string {
	return "deliver"
}
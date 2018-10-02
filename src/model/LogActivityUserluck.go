package model

import (
	"time"
)

type LogActivityUserluck struct {
	ID int64 `gorm:"column:id" json:"id"`
	MemberID int64 `gorm:"column:member_id" json:"member_id"`
	ActivityId int64 `gorm:"column:activity_id" json:"activity_id"`
	BuyNumber int64 `gorm:"column:buy_number" json:"buy_number"`
	BuyTime time.Time `gorm:"column:buy_time" json:"buy_time"`
	MerchantId int64 `gorm:"column:_merchant_id" json:"merchant_id"`
	StoreId int64 `gorm:"column:_store_id" json:"store_id"`
	Coin int64 `gorm:"column:coin" json:"coin"`
	Mcoin int64 `gorm:"column:mcoin" json:"mcoin"`
	FormId string `gorm:"column:form_id" json:"form_id"`
}
func (LogActivityUserluck) TableName() string {
	return "log_activity_userluck"
}
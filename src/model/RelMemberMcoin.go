package model

import "time"

type RelMemberMcoin struct {
	ID int64 `gorm:"column:id" json:"id"`
	MemberId int64 `gorm:"column:member_id" json:"member_id"`
	MerchantId int64 `gorm:"column:merchant_id" json:"merchant_id"`
	Mcoin int64 `gorm:"column:mcoin" json:"mcoin"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (RelMemberMcoin) TableName() string {
	return "rel_member_mcoin"
}
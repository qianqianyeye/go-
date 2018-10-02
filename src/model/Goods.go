package model

import "time"

type Goods struct {
	ID int64 `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Type int64 `gorm:"column:type" json:"type"`
	Size string `gorm:"column:size" json:"size"`
	CostPrice int64 `gorm:"column:cost_price" json:"cost_price"`
	SellingPrice int64 `gorm:"column:selling_price" json:"selling_price"`
	Headimgurl string `gorm:"column:headimgurl" json:"headimgurl"`
	Stock int64 `gorm:"column:stock" json:"stock"`
	UploadAccountId int64 `gorm:"column:upload_account_id" json:"upload_account_id"`
	Status int64 `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	MerchantId int64 `gorm:"column:merchant_id" json:"merchant_id"`

}

func (Goods) TableName() string {
	return "goods"
}
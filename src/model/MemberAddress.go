package model
type MemberAddress struct {
	ID int64 `gorm:"column:id" json:"id"`
	Province string `gorm:"column:province" json:"province"`
	City string `gorm:"column:city" json:"city"`
	Area string `gorm:"column:area" json:"area"`
	Body string `gorm:"column:body" json:"body"`
	Status int64 `gorm:"column:status" json:"status"`
	MemberId int64 `gorm:"column:member_id" json:"member_id"`

}
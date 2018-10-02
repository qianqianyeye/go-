package model

import "time"

type Activity struct {
	ID int64 `gorm:"column:id" json:"id"`
	ActivityName string `gorm:"column:activity_name" json:"activity_name" binding:"required"`
	ActivityStatus int64 `gorm:"column:activity_status" json:"activity_status" binding:"required"`
	MerchantId int64 `gorm:"column:merchant_id" json:"merchant_id"`
	StroreId int64 `gorm:"column:strore_id" json:"strore_id"`
	StartTime time.Time `gorm:"column:start_time" json:"start_time" form:"start_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	//EndTime time.Time `gorm:"column:end_time" json:"end_time" form:"end_time" binding:"required" time_format:"2006-01-02 15:04:05"`
	ActivityForwardNum int64 `gorm:"column:activity_forward_num" json:"activity_forward_num"`
	ParticipantsNumber int64 `gorm:"column:participants_number" json:"participants_number"`
	ActivityType int64 `gorm:"column:activity_type" json:"activity_type"`
	TotalMcoin int64 `gorm:"column:total_mcoin" json:"total_mcoin" binding:"required"`
	Remark string `gorm:"column:remark" json:"remark"`
	ImageUrl string `gorm:"column:image_url" json:"image_url" binding:"required"`
	Price int64 `gorm:"column:price" json:"price" binding:"required"`

	RelActivityShare RelActivityShare  `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ID;" json:"rel_activity_share,omitempty"`
	RelActivityThreshold []RelActivityThreshold `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ID;" json:"rel_activity_threshold,omitempty"`
	RelActivityRule []RelActivityRule `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ID;" json:"rel_activity_rule,omitempty"`
	ActivityAwardInfo []ActivityAwardInfo `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ID;" json:"activity_award_info,omitempty" binding:"required"`
	LogActivityUserluck []LogActivityUserluck `gorm:"ForeignKey:ActivityId;AssociationForeignKey:ID;" json:"log_activity_userluck,omitempty"`
}

func (Activity) TableName() string {
	return "activity"
}
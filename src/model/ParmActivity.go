package model

import "time"

type DeleteParm struct {
	//Data DeleteActivity `form:"data" json:"data" binding:"required"`
	Id int64 `form:"id" json:"id" binding:"required"`
	MerchantId int64 `form:"merchant_id" json:"merchant_id" binding:"required"`
	Status int `form:"status" json:"status"`
	Msg string `form:"msg" json:"msg"`
}
type DeleteActivity struct {
	Id int64 `form:"id" json:"id" binding:"required"`
	MerchantId int64 `form:"merchant_id" json:"merchant_id" binding:"required"`
}
type SubMitParm struct {
	//Data ParmActive `form:"data" json:"data" binding:"required"`
	Activity Activity `form:"activity"  binding:"required"`
	Status int `form:"status" json:"status"`
	Msg string `form:"msg" json:"msg"`
}
type ParmActive struct {
	Activity Activity `form:"activity"  binding:"required"`
	//Share RelActivityShare `form:"share"  binding:"required"`
	//Rule []RelActivityRule `form:"rule" binding:"required"`
	//Threshold []RelActivityThreshold `form:"threshold"`
}
type LuckParm struct {
	MemberID int64  `form:"member_id"  json:"member_id" ` //用户ID
	BuyNum int64  `form:"buy_num" json:"buy_num" binding:"required"`  //购买份数
	//Price int64 `form:"price" json:"price" binding:"required"` //单价
	FormId string  `form:"form_id" json:"form_id" binding:"required"`//推送formId
	ActivityId int64 `form:"activity_id" json:"activity_id" binding:"required"`//活动ID
	//ThreSholdIds []int64 `form:"activity_id"  binding:"required"`//门槛
}

type ParmValid struct {
	Current int `form:"current"  binding:"required"`
	Page_size int `form:"page_size"  binding:"required"`
	Start_time time.Time `form:"start_time" time_format:"2006-01-02 15:04:05"`
	End_time time.Time `form:"end_time"  time_format:"2006-01-02 15:04:05"`
	Merchant_id string `form:"merchant_id"`
	ID int64 `form:"id"`
	Store_id string `form:"store_id"`
	Sort string     `form:"sort"`
}
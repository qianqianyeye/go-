package service

import (
	"strconv"
	"SaasActivityService/src/db"
	"SaasActivityService/src/model"
	"SaasActivityService/src/webgo"
	"fmt"
)

type LuckService struct {}
type Count struct {
	Count int
}
func (ctrl *LuckService)CreateActivity(parmActive model.Activity) (bool,model.Activity){
	tx:=db.SqlDB.Begin()
	if err := tx.Create(&parmActive).Error; err != nil {
		tx.Rollback()
		return false,parmActive
	}
	tx.Commit()
	db.CfRedis[0].HSet("activity:"+strconv.FormatInt(parmActive.ID,10),"buyNum",0)
	//db.CfRedis[0].HSet("activity:"+strconv.FormatInt(parmActive.ID,10),"buyPeopleNum",0)
	db.CfRedis[0].HSet("activity:"+strconv.FormatInt(parmActive.ID,10),"forwardNum",0)
	db.CfRedis[0].HSet("activity:"+strconv.FormatInt(parmActive.ID,10),"sendCoin",0)
	db.CfRedis[0].HSet("activity:"+strconv.FormatInt(parmActive.ID,10),"clickNum",0)
	return true,parmActive
}

func (ctrl *LuckService)UpdateActivity(parmActive map[string]interface{},id interface{}) bool {
	tx:=db.SqlDB.Begin()
	if err := tx.Table("activity").Where("id=?",id).Update(parmActive).Error; err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}
func (ctrl *LuckService)DeleteActivity(parm int64)  bool{
	var activity model.Activity
	activity.ID=parm
	db.SqlDB.Model(&activity).Update("activity_status",-1)
	return true
}

func (ctrl *LuckService)GetAcitvityList(pageModel webgo.PageModel,id int64,merchantId int64) ([]model.Activity,webgo.PageModel){
	var activity []model.Activity

	GetDbPageData(pageModel).Order("activity_status,start_time desc").Preload("RelActivityThreshold.RelActivityThresholdinfo").
		Preload("RelActivityThreshold").Preload("RelActivityRule.RelActivityRuleinfo").Preload("RelActivityRule").
			Preload("RelActivityShare").Preload("ActivityAwardInfo").Find(&activity,"activity_status != ? and merchant_id=?",-1,merchantId)

	var Count Count
	db.SqlDB.Select("count(*)count").Table("activity").Where("activity_status != ?",-1).Scan(&Count)
	pageModel.Total=Count.Count
	return activity,pageModel
}

func  (ctrl *LuckService)SearchActivity(parm map[string]interface{}) model.Activity{
	cond, vals, err := whereBuild(parm)
	if err != nil {
		panic(err)
	}
	var activity model.Activity
	db.SqlDB.Preload("RelActivityThreshold.RelActivityThresholdinfo").
			Preload("RelActivityThreshold").Preload("RelActivityRule.RelActivityRuleinfo").Preload("RelActivityRule").
			Preload("RelActivityShare").Preload("ActivityAwardInfo.Goods").Preload("ActivityAwardInfo").Where(cond, vals...).Find(&activity)
		//db.SqlDBFind(&activity)
	return activity
}

func (ctrl *LuckService)GetActivityById(id int64)  model.Activity{
	var activity model.Activity
	db.SqlDB.Find(&activity,"id=?",id)
	return activity
}

func (ctrl *LuckService)GetJudgeInfo(id int64)  model.Activity{
	var activity model.Activity
	db.SqlDB.Preload("RelActivityThreshold.RelActivityThresholdinfo").
		Preload("RelActivityThreshold").Preload("ActivityAwardInfo").Find(&activity,"id=?",id)
	return activity
}

//用户在这个商家累计参与的活动次数
func (ctrl *LuckService)GetMemberActivityNum(merchantId int64,memberId int64)  int{
	//var num int
	var count Count
	db.SqlDB.Table("log_activity_userluck").Select("activity_id,count(member_id)count,member_id").
		Where("_merchant_id=? and member_id=? ",merchantId,memberId).Group("activity_id").Scan(&count)
	return count.Count
}

func (ctrl *LuckService)GetMemberBalance(memberId int64,merchantId int64)  model.Member {
	var Member model.Member
	db.SqlDB.Preload("RelMemberMcoin","merchant_id=?",merchantId).Find(&Member,"id=?",memberId)
	return Member
}
func (ctrl *LuckService)PayMoney(coin int64,mCoin int64,memberId int64,logUserLuck model.LogActivityUserluck,merchantId int64) bool {
	var relMemberMcoin model.RelMemberMcoin
	var member model.Member
	relMemberMcoin.Mcoin=mCoin
	member.Balance=coin
	tx:=db.SqlDB.Begin()
	if mCoin!=0 && coin==0{
		//商户币足够
		if err := tx.Table("rel_member_mcoin").Where("member_id=? and merchant_id=?",memberId,merchantId).Update(map[string]interface{}{"mcoin": mCoin}).Error; err != nil {
			tx.Rollback()
			return false
		}
	}else {
		//商户币不够
		if err := tx.Table("rel_member_mcoin").Where("member_id=? and merchant_id=?",memberId,merchantId).Update(map[string]interface{}{"mcoin": mCoin}).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return false
		}
		if err := tx.Table("member").Where("id=?",memberId).Update(map[string]interface{}{"balance": coin}).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return false
		}
	}
	if err := tx.Create(&logUserLuck).Error; err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func (ctrl *LuckService)GetBuyPeopleNum(activityId int64)  int{
	var Count Count
	db.SqlDB.Table("log_activity_userluck").Select("count(*)count,member_id").Where("activity_id=?",activityId).Group("member_id").Scan(&Count)
	return Count.Count
}
func (ctrl *LuckService)GetBuyPeople(activityId int64)  []model.LogActivityUserluck{
	var logActivity []model.LogActivityUserluck
	db.SqlDB.Table("log_activity_userluck").Select("member_id").Where("activity_id=?",activityId).Group("member_id").Scan(&logActivity)
	return logActivity
}

func (ctrl *LuckService)GetBuyNum(activityId int64)  int {
	var Count Count
	db.SqlDB.Table("log_activity_userluck").Select("sum(buy_number)count").Where("activity_id=?",activityId).Scan(&Count)
	return Count.Count
}

func (ctrl *LuckService)GetDrawList(merchantId int64,pageModel webgo.PageModel) ([]model.RelActivityOpenluck,webgo.PageModel ){
	var OpenLuck []model.RelActivityOpenluck
	var Count Count
	GetDbPageData(pageModel).Order("open_time desc").Preload("Activity").
		Preload("LogActivityUserluck").Preload("Member").Find(&OpenLuck,"_merchant_id=?",merchantId)
	db.SqlDB.Select("count(*)count").Table("rel_activity_openluck").Where("_merchant_id=?",merchantId).Scan(&Count)
	pageModel.Total=Count.Count
	return OpenLuck,pageModel
}

func (ctrl *LuckService)GetDrawByActivity(merchantId int64,activityId int64) []model.RelActivityOpenluck{
	var OpenLuck []model.RelActivityOpenluck

	db.SqlDB.Order("open_time desc").Preload("Activity").
		Preload("LogActivityUserluck","activity_id=?",activityId).Preload("Member").Find(&OpenLuck,"_merchant_id=? and activity_id=?",merchantId,activityId)

	for i,val := range OpenLuck{
		var RelActivityDeliver model.RelActivityDeliver
		db.SqlDB.Preload("Deliver","member_id=?",val.MemberId).Find(&RelActivityDeliver,"activity_id=?",activityId)
		if RelActivityDeliver.ID!=0 {
			data := webgo.StructToJsonMap(RelActivityDeliver.Deliver)
			OpenLuck[i].Deliver=data
		}
	}
	//db.SqlDB.Select("count(*)count").Table("rel_activity_openluck").Where("_merchant_id=? and activity_id=?",merchantId,activityId).Scan(&Count)
	return OpenLuck
}

func (ctrl *LuckService)GetLuckAdress(activityId int64,merchantId int64)  []model.ActivityAwardInfo{
	var ActivityAwards  []model.ActivityAwardInfo
	db.SqlDB.Preload("RelActivityDeliver.Deliver.Member").Preload("RelActivityDeliver.Deliver","merchant_id=?",merchantId).
		Preload("RelActivityDeliver").Find(&ActivityAwards,"activity_id=?",activityId)
	return ActivityAwards
}
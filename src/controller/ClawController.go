package controller

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
	"strconv"
	"math/rand"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"SaasActivityService/src/webgo"
	"SaasActivityService/src/model"
	"SaasActivityService/src/db"
	"SaasActivityService/src/middleware"
	"SaasActivityService/src/service"
)

type ClawController struct {
	webgo.Controller
}
type ClawParm struct {
	ActivityId int64 `form:"activity_id" json:"activity_id" binding:"required"`
}

type ActivityInfo struct {
	ActivityId int64
	//TotalShareNumber int64
	Rule []map[int64]string
	CurShareNumber  int64
	CurBuyNumber int64
	MerchantId int64
	StoreId int64
	Status int64
	StartTime time.Time
}
var AtInfo *sync.Map

func init()  {
	AtInfo=new(sync.Map)
}
func InitAtInfo()  {
	var activity []model.Activity
	db.SqlDB.Preload("RelActivityThreshold.RelActivityThresholdinfo").
		Preload("RelActivityThreshold").Preload("RelActivityRule.RelActivityRuleinfo").Preload("RelActivityRule").
		Preload("RelActivityShare").Preload("LogActivityUserluck").Find(&activity,"activity_status in (0,1)")
	for _,val :=range activity {
		AddActivity(val)
	}
	go OpenLuck()//开奖线程
}
func AddActivity(val model.Activity)  {
	var CurBuyNumber int64=0
	var tempInfo ActivityInfo
	tempInfo.ActivityId=val.ID
	tempInfo.StoreId=val.StroreId
	tempInfo.MerchantId=val.MerchantId
	if db.CfRedis[0].Exists("activity:"+strconv.FormatInt(val.ID,10)).Val()!=0 {
		forwardNumstr:=db.CfRedis[0].HGet("activity:"+strconv.FormatInt(val.ID,10),"forwardNum").Val()
		tempInfo.CurShareNumber,_=strconv.ParseInt(forwardNumstr, 10, 64)  //修改为从Redis拿
	}
	tempInfo.Status=val.ActivityStatus
	tempInfo.StartTime=val.StartTime
	var rule =make([]map[int64]string,0)

	for _,val := range val.RelActivityRule {
			temp := make(map[int64]string)
			temp[val.RuleId]=val.Body
			rule=append(rule, temp)
		}

		for _,lval:=range val.LogActivityUserluck{
			CurBuyNumber=lval.BuyNumber+CurBuyNumber
		}

	tempInfo.CurBuyNumber=CurBuyNumber
	tempInfo.Rule=rule
	AtInfo.Store(val.ID,tempInfo)
}
func (ctrl *ClawController) Router(router *gin.Engine) {
	r := router.Group("go/api/v1/claw",middleware.ClawMiddle)
	r.POST("activityshare",ctrl.activityShare)
}

func (ctrl *ClawController)activityShare(ctx *gin.Context)  {
	ActivityId, err := strconv.ParseInt(webgo.GetResult(ctx.Keys["activity_id"]), 10, 64)
	if err!=nil {
		fmt.Println(err)
		webgo.Result(ctx,webgo.SUCCESS,"参数错误！",nil,nil)
	}
	if val,ok := AtInfo.Load(ActivityId);ok {
			var tempInfo ActivityInfo
			switch val.(type) {
			case ActivityInfo:
				tempInfo=val.(ActivityInfo)
			}
			tempInfo.CurShareNumber +=1
			fmt.Println(tempInfo.CurShareNumber)
			AtInfo.Store(tempInfo.ActivityId,tempInfo)
			}
}

func OpenLuck()  {
	for  {
		//Range
		//遍历sync.Map, 要求输入一个func作为参数
		f := func(k, v interface{}) bool {
			//这个函数的入参、出参的类型都已经固定，不能修改
			//可以在函数体内编写自己的代码，调用map中的k,v
			var temp ActivityInfo
			switch v.(type) {
			case ActivityInfo:
				temp=v.(ActivityInfo)
			}
			//status := temp.Status
			//curt:=time.Now().Unix()
			//			//start:=temp.StartTime.Unix()
			//			//if curt>=start {
			//			//	parmActivity:= make(map[string]interface{})
			//			//	parmActivity["activity_status"]=1
			//			//	luckService.UpdateActivity(parmActivity,temp.ActivityId)
			//			//}
				//所有规则都满足则所有值为true
				var flagMap = make(map[int64]bool)
				//遍历规则
				for _,ruv :=range temp.Rule{
					for rk,rv:=range ruv{
						//rk:关系 rv:条件
						if rk==1 {
							t:=time.Now().Unix()
							old :=webgo.StringToTime(rv).Unix()
							//当前时间大于等于规定时间，开奖条件1满足
							if t>=old {
								flagMap[rk]=true
							}else {
								flagMap[rk]=false
							}
						}
						if rk==2 {
							total, err := strconv.ParseInt(rv, 10, 64)
							if err!=nil {
								webgo.Error(err)
							}
							//当前转发数大于规定数量，开奖条件2满足
							if temp.CurShareNumber>= total{
								flagMap[rk]=true
							}else{
								flagMap[rk]=false
							}
						}
						if rk==3 {
							buyTotal, err := strconv.ParseInt(rv, 10, 64)
							if err!=nil {
								webgo.Error(err)
							}
							//当前购买份数大于等于规定总份数，开奖条件3满足
							if temp.CurBuyNumber>= buyTotal{
								flagMap[rk]=true
							}else{
								flagMap[rk]=false
							}
						}
					}
				}
			//开奖时间满足，其他条件不满足 退币
			timeflag:=false
			if _, ok := flagMap[1]; ok {
				//存在
				timeflag=flagMap[1]
			}
			if timeflag {
				for _,v:=range flagMap{
					if v==false {
						//退币
						BackBalance(temp.ActivityId)
						parmActivity:= make(map[string]interface{})
						parmActivity["activity_status"]=10
						//将活动更改为结束状态
						luckService.UpdateActivity(parmActivity,temp.ActivityId)
						//结束的活动从map中移除
						AtInfo.Delete(temp.ActivityId)
						FormClaw(temp.ActivityId)
						return true
					}
				}
			}
			//判断是否满足
			for _,v:=range flagMap{
				if v==false {
					return true
				}
			}
			UserLuck(temp) //抽奖
			//通知开奖
			FormClaw(temp.ActivityId)
			return true
		}
		AtInfo.Range(f)
		time.Sleep(time.Second)
	}
}

func FormClaw(activityId int64)  {
	parm := make(map[string]interface{})
	parm["activity_id"]=activityId
	url :=webgo.ClawIp+"/luckdraw/activity/przie"
	HttpPost(url,parm)
}

func BackBalance(activityId int64) bool{
	var logUserLuck []model.LogActivityUserluck
	db.SqlDB.Find(&logUserLuck,"activity_id=?",activityId)
	tx:=db.SqlDB.Begin()
	for _,val := range logUserLuck{
		id:=val.MemberID
		if err :=tx.Exec("UPDATE `member` SET balance= balance + ?  WHERE id=? ",val.Coin,id).Error;err!=nil{
			tx.Rollback()
			return false
		}
		if err :=tx.Exec("UPDATE `rel_member_mcoin` SET mcoin= mcoin + ?  WHERE member_id=? ",val.Mcoin,id).Error;err!=nil{
			tx.Rollback()
			return false
		}
		if err :=tx.Exec("UPDATE `log_activity_userluck` SET status=0  where id=?",val.ID).Error;err!=nil{
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}
func UserLuck(info ActivityInfo) {
	var luckService service.LuckService
	var logActivityUserLuck []model.LogActivityUserluck
	//var openLuck model.RelActivityOpenluck
	db.SqlDB.Table("log_activity_userluck").Select("sum(buy_number)buy_number,member_id").Where("activity_id=?",info.ActivityId).Group("member_id").Scan(&logActivityUserLuck)
	if len(logActivityUserLuck)>0 {
		users := make(map[int64]int64)
		for _,v:=range logActivityUserLuck{
			users[v.MemberID]=v.BuyNumber
		}
		generator := GetAwardUserName(users)
		MemberId :=generator()
		var openLuck model.RelActivityOpenluck
		openLuck.MerchantId=info.MerchantId
		openLuck.StoreId=info.StoreId
		openLuck.MemberId=MemberId
		openLuck.OpenTime=time.Now()
		openLuck.ActivityId=info.ActivityId
		//将中奖人存入
		db.SqlDB.Exec("insert into rel_activity_openluck (_merchant_id,_store_id,member_id,open_time,activity_id) values(?,?,?,?,?)",openLuck.MerchantId,openLuck.StoreId,openLuck.MemberId,openLuck.OpenTime,openLuck.ActivityId)
	}
	parmActivity:= make(map[string]interface{})
	if len(logActivityUserLuck)<=0 {
		parmActivity["activity_status"]=10
	}else {
		parmActivity["activity_status"]=9
	}
	//将活动更改为结束状态
	luckService.UpdateActivity(parmActivity,info.ActivityId)
	//结束的活动从map中移除
	AtInfo.Delete(info.ActivityId)
}


func GetAwardUserName(users map[int64]int64) (generator func()int64){
	var sum_num int64
	name_arr := make([]int64,len(users))
	offSetOf := make([]int64,len(users))
	index :=0
	for k,v:=range users{
		offSetOf[index]=sum_num
		sum_num+=v
		name_arr[index]=k
		index+=1
	}
	generator =func()int64{
		awardNum := rand.Int63n(sum_num) //获得中奖号码
		return name_arr[binarySearch(offSetOf,awardNum)] //在数轴中查找该号码对应的用户
	}
	return
}

func binarySearch(nums []int64,target int64) int {
	start,end :=0,len(nums)-1
	for start<=end  {
		mid := (start+end)/2
		//mid :=  start + (end-start)/2
		if nums[mid]>target {
			end =mid-1
		}else if nums[mid]<target {
			if mid+1==len(nums) {
				return mid
			}
			if nums[mid+1]>target {
				return mid
			}
			start=mid+1
		}else {
			return mid
		}
	}
	return -1
}


func HttpPost(url string, payload map[string]interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	pMap:=middleware.GetClawSign(payload)
	payload["timeStamp"]=pMap["timeStamp"]
	payload["sign"]=pMap["sign"]
	bytesData, err := json.Marshal(payload)
	if err != nil {
		return m, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return m, err
		}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return m, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m, err
		}
	json.Unmarshal(body, &m)
	return m, nil
}
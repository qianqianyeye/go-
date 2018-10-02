package controller

import (
	"SaasActivityService/src/db"
	"SaasActivityService/src/middleware"
	"SaasActivityService/src/model"
	"SaasActivityService/src/service"
	"SaasActivityService/src/webgo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"time"
)

var luckService service.LuckService

type LuckController struct {
	webgo.Controller
}

func (ctrl *LuckController) Router(router *gin.Engine) {
	r := router.Group("activity/v1/luck", middleware.Middle)
	//r := router.Group("go/v1/luck")
	r.POST("activity", ctrl.createActive)
	r.PATCH("activity", ctrl.updateActivity)
	r.DELETE("activity/:id", ctrl.deleteActivity)
	r.GET("activity", ctrl.getActivityList)
	r.GET("activity/:id", ctrl.searchActivity)
	r.GET("activity/:id/statement", ctrl.statement)
	r.GET("activity/:id/address", ctrl.address)
	r.GET("activity/:id/luck", ctrl.draw) //中奖列表

	rs := router.Group("activity/v1", middleware.CMiddle)
	rs.POST("activity/luck", ctrl.luck)
}
func (ctrl *LuckController) draw(ctx *gin.Context) {
	stringId := webgo.GetResult(ctx.Keys["merchantId"])
	merchantId, _ := strconv.ParseInt(stringId, 10, 64)
	//pageModel := webgo.GetPageInfo(strconv.Itoa(parmValid.Page_size), strconv.Itoa(parmValid.Current))
	activityId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		webgo.Result(ctx, webgo.SUCCESS, "参数错误！", nil, nil)
		return
	}
	result := luckService.GetDrawByActivity(merchantId, activityId)
	//activity, page := luckService.GetDrawList(merchantId,pageModel)
	webgo.Result(ctx, webgo.SUCCESS, nil, result, nil)

}

func (ctrl *LuckController) address(ctx *gin.Context) {
	id := ctx.Param("id")
	if id != "" {
		//pageModel := webgo.GetPageInfo(strconv.Itoa(parmValid.Page_size), strconv.Itoa(parmValid.Current))
		activityId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		stringId := webgo.GetResult(ctx.Keys["merchantId"])
		merchantId, _ := strconv.ParseInt(stringId, 10, 64)
		if err != nil {
			webgo.Result(ctx, webgo.SUCCESS, "参数错误！", nil, nil)
			return
		}
		activityAwards :=luckService.GetLuckAdress(activityId,merchantId)
		Delivers:=make([]model.Deliver,0)
		for _,val :=range activityAwards{
			for _,rval :=range val.RelActivityDeliver{
				Delivers=append(Delivers,rval.Deliver)
			}
		}
		webgo.Result(ctx, webgo.SUCCESS, "成功！", Delivers, nil)
	} else {
		webgo.Result(ctx, webgo.SUCCESS, "参数错误！", nil, nil)
	}
}

func (ctrl *LuckController) luck(ctx *gin.Context) {
	var luckParm model.LuckParm
	var member model.Member
	var payMoney int64 //需要付多少币
	var logUserLuck model.LogActivityUserluck
	var num int
	stringId := webgo.GetResult(ctx.Keys["memberId"])
	memberId, _ := strconv.ParseInt(stringId, 10, 64)
	luckParm.MemberID = memberId
	if err := ctx.ShouldBind(&luckParm); err == nil {
		activity := luckService.GetJudgeInfo(luckParm.ActivityId)
		if activity.ActivityStatus == 1 {
			flag := true
			flagMap := make(map[int64]bool)
			//验证门槛
			if len(activity.RelActivityThreshold) > 0 {
				for _, val := range activity.RelActivityThreshold {
					//累计参与次数
					if val.ThresholdId == 1 {
						num = luckService.GetMemberActivityNum(activity.MerchantId, luckParm.MemberID)
						bodyNum, _ := strconv.Atoi(val.Body)
						//累计参与次数 小于规定次数
						if num < bodyNum {
							flagMap[val.ThresholdId] = false
						} else {
							flagMap[val.ThresholdId] = true
						}
					}
					//用户余额满多少
					if val.ThresholdId == 2 {
						member = luckService.GetMemberBalance(luckParm.MemberID,activity.MerchantId)
						memberBalance := member.Balance + member.RelMemberMcoin.Mcoin
						balance, err := strconv.ParseInt(val.Body, 10, 64)
						if err != nil {
							panic(err)
						}
						//用户余额小于规定金额
						if memberBalance < balance {
							flagMap[val.ThresholdId] = false
						}
					}
				}
			}
			for k, val := range flagMap {
				if val == false && k == 1 {
					flag = false
					webgo.Result(ctx, webgo.JOINNUM, "累计参与次数不满足条件！", num, nil)
					return
				}
				if val == false && k == 2 {
					flag = false
					webgo.Result(ctx, webgo.TBALANCE, "用户余额小于规定门槛金额！", member, nil)
					return
				}
			}
			//是否满足门槛，满足flag=true
			if flag {
				var payMcoin int64 //付完后剩余的商户币
				var payCoin int64  //付完后剩余的平台币
				//用户余额是否足够
				payMoney = luckParm.BuyNum * activity.Price
				if member.ID == 0 {
					member = luckService.GetMemberBalance(luckParm.MemberID,activity.MerchantId)
				}
				memberBalance := member.Balance + member.RelMemberMcoin.Mcoin
				fmt.Println("payMoney%d",payMoney)
				fmt.Println("memberBalance%d",memberBalance)
				if payMoney > memberBalance {
					fmt.Println("payMoney%d",payMoney)
					fmt.Println("memberBalance%d",memberBalance)
					webgo.Result(ctx, webgo.BALANCE, "用户余额不足！", member, nil)
					fmt.Println("已返回！")
					return
				}
				//fmt.Println()
				//fmt.Println("Activity%v",activity)
				//剩余份数是否足够
				//for _, val := range activity.RelActivityRule {
				//	fmt.Println("rule%v",val)
				//	if val.RuleId == 3 {
				//		curBuyNumStr := db.CfRedis[0].HGet("activity:"+strconv.FormatInt(luckParm.ActivityId, 10), "buyNum").Val()
				//		curBuyNum, _ := strconv.ParseInt(curBuyNumStr, 10, 64)
				//		total, _ := strconv.ParseInt(val.Body, 10, 64)
				//		fmt.Println("total%d",total)
				//		fmt.Println("curBuyNum%d",curBuyNum)
				//		if total-curBuyNum < 0 {
				//			fmt.Println("total%d",total)
				//			fmt.Println("curBuyNum%d",curBuyNum)
				//			webgo.Result(ctx, webgo.SNUM, "剩余份数不足", total-curBuyNum, nil)
				//			return
				//		}
				//	}
				//}
				if member.RelMemberMcoin.Mcoin >= payMoney {
					logUserLuck.Mcoin = payMoney
					payMcoin = member.RelMemberMcoin.Mcoin - payMoney
				} else {
					//商户币不够，扣完为0
					logUserLuck.Mcoin = member.RelMemberMcoin.Mcoin //扣除的商户币
					tempPay := payMoney - member.RelMemberMcoin.Mcoin
					logUserLuck.Coin = tempPay
					payCoin = member.Balance - tempPay
				}
				logUserLuck.MemberID = luckParm.MemberID
				logUserLuck.ActivityId = luckParm.ActivityId
				logUserLuck.MerchantId = activity.MerchantId
				logUserLuck.BuyNumber = luckParm.BuyNum
				logUserLuck.FormId = luckParm.FormId
				logUserLuck.BuyTime = time.Now()
				f := luckService.PayMoney(payCoin, payMcoin, luckParm.MemberID, logUserLuck,activity.MerchantId)
				if val, ok := AtInfo.Load(activity.ID); ok {
					var tempInfo ActivityInfo
					switch val.(type) {
					case ActivityInfo:
						tempInfo = val.(ActivityInfo)
					}
					tempInfo.CurBuyNumber += luckParm.BuyNum
					AtInfo.Store(tempInfo.ActivityId, tempInfo)
				}
				db.CfRedis[0].HIncrBy("activity:"+strconv.FormatInt(luckParm.ActivityId, 10), "buyNum", luckParm.BuyNum)
				db.CfRedis[0].HIncrBy("activity:"+strconv.FormatInt(luckParm.ActivityId, 10),"member:"+strconv.FormatInt(memberId, 10),luckParm.BuyNum)
				//buyPeopleNum:=luckService.GetBuyPeopleNum(luckParm.ActivityId)
				//db.CfRedis[0].HSet("activity:"+strconv.FormatInt(luckParm.ActivityId,10),"buyPeopleNum",buyPeopleNum)
				db.CfRedis[0].SAdd("activity:"+strconv.FormatInt(luckParm.ActivityId, 10)+":buyPeople", memberId)
				if f {
					webgo.Result(ctx, webgo.SUCCESS, "成功！", nil, nil)
				} else {
					webgo.Result(ctx, webgo.ERROR, "失败！", nil, nil)
				}
			}
		} else {
			webgo.Result(ctx, webgo.ERROR, "活动已关闭或结束！", nil, nil)
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}

func (ctrl *LuckController) statement(ctx *gin.Context) {
	var stateMent model.Statement
	id := ctx.Param("id")
	if id != "" {
		if db.CfRedis[0].Exists("activity:"+id).Val() != 0 {
			stateMent.BuyNum, _ = strconv.Atoi(db.CfRedis[0].HGet("activity:"+id, "buyNum").Val())
			//stateMent.BuyPeopleNum,_=strconv.Atoi(db.CfRedis[0].HGet("activity:"+id,"buyPeopleNum").Val())
			stateMent.ForwardNum, _ = strconv.Atoi(db.CfRedis[0].HGet("activity:"+id, "forwardNum").Val())
			stateMent.SendCoin, _ = strconv.Atoi(db.CfRedis[0].HGet("activity:"+id, "sendCoin").Val())
			stateMent.ClickNum, _ = strconv.Atoi(db.CfRedis[0].HGet("activity:"+id, "clickNum").Val())
			stateMent.BuyPeopleNum = db.CfRedis[0].SCard("activity:" + id + ":buyPeople").Val()
		} else {
			activityId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			logActivity := luckService.GetBuyPeople(activityId)
			for _, val := range logActivity {
				db.CfRedis[0].SAdd("activity:"+strconv.FormatInt(activityId, 10)+":buyPeople", val.MemberID)
			}
			buyNum := luckService.GetBuyNum(activityId)
			db.CfRedis[0].HSet("activity:"+strconv.FormatInt(activityId, 10), "buyNum", buyNum)
			stateMent.BuyNum = buyNum
			stateMent.BuyPeopleNum =  db.CfRedis[0].SCard("activity:" + id + ":buyPeople").Val()
		}
		webgo.Result(ctx, webgo.SUCCESS, nil, stateMent, nil)
	} else {
		webgo.Result(ctx, webgo.SUCCESS, "参数错误！", nil, nil)
	}

}
func (ctrl *LuckController) createActive(ctx *gin.Context) {
	subMitParm := model.SubMitParm{}
	stringId := webgo.GetResult(ctx.Keys["merchantId"])
	merchantId, _ := strconv.ParseInt(stringId, 10, 64)
	//var merchantId int64 = 52
	if err := ctx.ShouldBind(&subMitParm); err == nil {
		if merchantId == subMitParm.Activity.MerchantId {
			flag, activity := luckService.CreateActivity(subMitParm.Activity)
			if flag {
				AddActivity(activity)
				webgo.Result(ctx, webgo.SUCCESS, "创建活动成功", nil, nil)
			} else {
				webgo.Result(ctx, webgo.ERROR, "创建活动失败", nil, nil)
			}
		} else {
			webgo.Result(ctx, webgo.ERROR, "无权限！", nil, nil)
		}
	} else {
		webgo.Error(err.Error())
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}

func (ctrl *LuckController) updateActivity(ctx *gin.Context) {
	parm := make(map[string]interface{})
	var jwtMerchantId int64 = 52
	var id interface{}
	flag := false
	mflag := false
	buf, _ := ctx.GetRawData()
	var str string = string(buf[0:len(buf)])
	var rmap map[string]interface{}
	if err := json.Unmarshal([]byte(str), &rmap); err == nil {
	} else {
		fmt.Println(err)
	}
	for k, v := range rmap {
		if k == "id" {
			flag = true
			id = v
			continue
		}
		if k == "merchant_id" {
			merchantId, _ := strconv.ParseInt(webgo.GetResult(v), 10, 64)
			if merchantId == jwtMerchantId {
				mflag = true
			}
			continue
		}
		parm[k] = v
	}

	if flag == true && mflag == true {
		f := luckService.UpdateActivity(parm, id)
		if f {
			webgo.Result(ctx, webgo.SUCCESS, "操作成功", nil, nil)
		} else {
			webgo.Result(ctx, webgo.ERROR, "操作失败", nil, nil)
		}
	} else {
		webgo.Result(ctx, webgo.ERROR, "必须传递ID和MerchantId,或者MerchantId校验不通过", nil, nil)
	}
}

func (ctrl *LuckController) deleteActivity(ctx *gin.Context) {
	//var DeleteParm model.DeleteParm
	stringId := webgo.GetResult(ctx.Keys["merchantId"])
	merchantId, _ := strconv.ParseInt(stringId, 10, 64)
	id := ctx.Param("id")
	if id != "" {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			panic(err)
		}
		activity := luckService.GetActivityById(i)
		if activity.MerchantId == merchantId {
			flag := luckService.DeleteActivity(i)
			if flag {
				if _, ok := AtInfo.Load(activity.ID); ok {
					AtInfo.Delete(activity.ID)
				}
				webgo.Result(ctx, webgo.SUCCESS, "操作成功", nil, nil)
			} else {
				webgo.Result(ctx, webgo.ERROR, "操作失败", nil, nil)
			}
		} else {
			webgo.Result(ctx, webgo.ERROR, "参数错误", nil, nil)
		}
	} else {
		webgo.Result(ctx, webgo.ERROR, "参数错误", nil, nil)
	}

}
func (ctrl *LuckController) getActivityList(ctx *gin.Context) {
	var parmValid model.ParmValid
	stringId := webgo.GetResult(ctx.Keys["merchantId"])
	merchantId, _ := strconv.ParseInt(stringId, 10, 64)
	if err := ctx.ShouldBindWith(&parmValid, binding.Query); err == nil {
		pageModel := webgo.GetPageInfo(strconv.Itoa(parmValid.Page_size), strconv.Itoa(parmValid.Current))
		activity, page := luckService.GetAcitvityList(pageModel, parmValid.ID, merchantId)
		webgo.Result(ctx, webgo.SUCCESS, nil, activity, page)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}
}

func (ctrl *LuckController) searchActivity(ctx *gin.Context) {
	parm := make(map[string]interface{})
	id := ctx.Param("id")
	if id != "" {
		parm["id"] = id
		stringId := webgo.GetResult(ctx.Keys["merchantId"])
		merchantId, _ := strconv.ParseInt(stringId, 10, 64)
		parm["merchant_id"] = merchantId
		result := luckService.SearchActivity(parm)
		webgo.Result(ctx, webgo.SUCCESS, "成功!", result, "")
	} else {
		webgo.Result(ctx, webgo.SUCCESS, "找不到ID", "", "")
	}
}

package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
	"strconv"
	"SaasActivityService/src/webgo"
)

//验证
func ClawMiddle(c *gin.Context) {
	defer webgo.TryCatchWeb(c)
	buf, _ := c.GetRawData()
	var str string = string(buf[0:len(buf)])
	var rmap map[string]interface{}
	if err := json.Unmarshal([]byte(str), &rmap); err == nil {
	} else {
		fmt.Println(err)
	}
	webgo.Debug("接收到参数：%+v\n",rmap)
	var keys []string
	for k := range rmap {
		if k != "sign" && k != "timeStamp" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var result string
	for _, val := range keys {
		tempStr := webgo.GetResult(rmap[val])
		c.Set(val, rmap[val])
		result = result + "&" + val + "=" + tempStr
	}
	result = string([]rune(result)[1:]) + webgo.GetResult(rmap["timeStamp"])
	webgo.Debug("加密前字符串：%+v\n",result)
	serverSign := webgo.GetMd5(result)
	webgo.Debug("加密后字符串：%+v\n",serverSign)
	webgo.Debug("要对比的加密字符串：%+v\n",webgo.GetResult(rmap["sign"]))
	if serverSign != webgo.GetResult(rmap["sign"]) {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"data": "签名错误!",
			"status":1,
		})
		return
	}
	c.Set("activity_id",rmap["activity_id"])
	c.Next()
}

func GetClawSign(parm map[string]interface{}) map[string]interface{} {
	var keys []string
	for k := range parm {
		if k != "sign" && k != "timeStamp" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var result string
	for _, val := range keys {
		tempStr := webgo.GetResult(parm[val])
		result = result + "&" + val + "=" + tempStr
	}
	resultMap :=make(map[string]interface{})
	t:=time.Now().Unix()
	resultMap["timeStamp"]=t
	result = string([]rune(result)[1:]) + strconv.FormatInt(t,10)
	sign :=webgo.GetMd5(result)
	resultMap["sign"]=sign
	return resultMap
}
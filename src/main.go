package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/itsjamie/gin-cors"
	"flag"
	l4g "github.com/alecthomas/log4go"
	"SaasActivityService/src/controller"
	"SaasActivityService/src/webgo"
	"SaasActivityService/src/db"
)

func registerRouter(router *gin.Engine) {
	new(controller.LuckController).Router(router)
	new(controller.ClawController).Router(router)
}
func main() {
	defer webgo.TryCatch()
	l4g.LoadConfiguration("config/log4g.xml") //使用加载配置文件,类似与java的log4j.propertites
	defer l4g.Close()               //注:如果不是一直运行的程序,请加上这句话,否则主线程结束后,也不会输出和log到日志文件
	dataBase := flag.Bool("MySql",false,"true :线上，false: 线下 默认:false")
	flag.Parse()
	db.InitDB(*dataBase) //初始化数据库
	db.InitRedis(*dataBase) //初始化Redis
	defer db.SqlDB.Close()
	go controller.InitAtInfo() //获取未开始跟开始的活动数据
	router := gin.Default()
	//网页跨域问题
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,PUT,POST,DELETE,OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	registerRouter(router)
	router.Run(":80")
}
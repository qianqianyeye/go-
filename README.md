


## 启动方式
#### 1.godep包管理工具安装

	go get -u -v github.com/tools/godep


#### 2.本地启动方式
       cd src/
       godep go run main.go  -MySql=true (true:线上，false：线下)

## 构建项目以及选用线上线下数据库
    cd src/
	godep go build main.go
	线上数据库
	main -MySql=true
	线下数据库
	main -MySql=false
	默认使用线下数据库

#### 本地部署成功后访问测试
    192.168.191.162:8031/go/v1/luck/activity












package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"bkd/middlewares"
	"bkd/route"
)

func main() {
	//通过godotenv管理环境，环境存储在.env文件，load函数默认读取.env文件,之后通过os.Getenv函数读取环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	//生成gin路由器
	router := gin.Default()
	//使用中间件
	middlewares.InitMiddlewares(router)
	//注册路由
	route.Register(router)
	//golang推荐设置信任的代理，此处设置为前端的地址
	router.SetTrustedProxies([]string{os.Getenv("CLIENT_ADDR")})
	//运行
	router.Run(os.Getenv("SERVER_ADDR"))
}

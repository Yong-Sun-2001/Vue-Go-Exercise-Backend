package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"

	"bkd/middlewares/my_session"
)

//注册my_session为全局的中间件
func Activate_My_Session(router *gin.Engine) {
	my_session.InitMgr("redis", os.Getenv("REDIS_ADDR"))
	router.Use(my_session.SessionMiddleware(my_session.MgrObj))
	router.Use(my_session.AuthCurrentUser())
	//my_session.AuthMiddleware()局部中间件，在route中设置
}


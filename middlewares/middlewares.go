package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"

	"bkd/middlewares/gin_session"
)

//注册session为全局的中间件
func InitMiddlewares(router *gin.Engine) {
	gin_session.InitMgr("redis", os.Getenv("REDIS_ADDR"))
	// gin_session.InitMgr("memory", "")  //内存版session，后面的地址参数无用处
	router.Use(gin_session.SessionMiddleware(gin_session.MgrObj))
}

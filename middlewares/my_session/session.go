package my_session

import (
	"bkd/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName  = "session_id" // sesion_id在Cookie中对应的key
	SessionContextName = "session"    // session data在gin上下文中对应的key
)

//定义一个全局的Mgr
var (
	// MgrObj 全局的Session管理对象（大仓库）
	MgrObj Mgr
)

//构造一个Mgr
func InitMgr(name string, addr string, option ...string) {

	switch name {
	case "memory": //初始化一个内存版管理者
		MgrObj = NewMemory()
	case "redis":
		MgrObj = NewRedisMgr()
	}
	MgrObj.Init(addr, option...) //初始化mgr
}

//
type SessionData interface {
	GetID() string // 返回自己的ID
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
	Del(key string)
	Clear()
	Save() // 保存
}

//不同版本的管理者接口
type Mgr interface {
	Init(addr string, option ...string)
	GetSessionData(sessionId string) (sd SessionData, err error)
	CreatSession() (sd SessionData)
}

//全局中间件
// 每次请求都会获取cookie中的session id 字段并取出session data，若出错则创建新的session,最后设置上下文的session字段值并写回session id字段值到cookie中
func SessionMiddleware(mgrObj Mgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.请求刚过来，从请求的cookie中获取SessionId
		SessionID, err := c.Cookie(SessionCookieName)
		var sd SessionData
		if err != nil {
			//1 第一次来，没有sessionid，-->给用户建一个sessiondata，分配一个sessionid
			sd = mgrObj.CreatSession()
		} else {
			//2. 根据sessionid去大仓库取sessiondata
			sd, err = mgrObj.GetSessionData(SessionID)
			if err != nil {
				//sessionid有误，取不到sessiondata，可能是自己伪造的
				//重新创建一个sessiondata
				sd = mgrObj.CreatSession()
				//更新sessionid
				SessionID = sd.GetID() //这个sessionid用于回写coookie
			}
		}
		//3.利用gin框架的c.Set("session",sessiondata)在ctx中设置session字段为sessionData
		c.Set(SessionContextName, sd)
		//回写cookie,设置session id字段为sessionData的id
		c.SetCookie(SessionCookieName, SessionID, 3600, "/", "localhost", false, false)
		c.Next()
	}
}

//全局中间件
// 获取登录用户并设置在ctx中，每次请求都会从session中取出uid验证登录状态
func AuthCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 从上下文中取到session data
		tmpSD, _ := ctx.Get(SessionContextName)
		sd := tmpSD.(SessionData)
		// 2. 从session data取到user_id
		value, _ := sd.Get("uid")
		if value != nil {
			user, err := model.GetUser("uid")
			if err == nil {
				ctx.Set("user", &user)
				fmt.Print("user set!!\n")
			}
		}
		ctx.Next()
	}
}

//局部中间件
// 要求登录的中间件,需要从ctx中获取user数据才算登录
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, _ := ctx.Get("user"); user != nil {
			if _, ok := user.(*model.User); ok {
				ctx.Next()
				return
			}
		}
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "局部中间件提示：需要登录",
		})
		ctx.Abort()
	}
}

// CurrentUser 获取当前用户  基于上述在ctx中设置user的方法
func CurrentUser(c *gin.Context) *model.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}

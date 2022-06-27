package route

import (
	"bkd/middlewares/gin_session"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//session校验是否登录
func AuthMiddleware(ctx *gin.Context) (gin_session.SessionData, bool) {
	// 1. 从上下文中取到session data
	tmpSD, _ := ctx.Get(gin_session.SessionContextName)
	sd := tmpSD.(gin_session.SessionData)
	// 2. 从session data取到isLogin
	value, err := sd.Get("isLogin")
	if err != nil { // 取不到就是没有登录
		return sd, false
	}
	isLogin, ok := value.(bool) //类型断言
	if !ok || !isLogin {        //取到了，但是值不是布尔类型或者值为false
		return sd, false
	}
	return sd, true
}

type LoginRouter struct{}

func (c *LoginRouter) Router(engine *gin.Engine) {
	engine.POST("/login", c.post_login)
	engine.GET("/home/:uid", c.get_home)
	engine.GET("/space", c.get_space)
	engine.GET("/logout", c.get_logout)
}

func (c *LoginRouter) post_login(ctx *gin.Context) {

	type info struct {
		Usr string `json:"usr"`
		Pwd string `json:"pwd"`
	}
	now := info{}
	//绑定，并解析参数
	err := ctx.ShouldBindJSON(&now)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	//验证输入的账号密码受否正确
	if now.Usr == "sun" && now.Pwd == "//I12vM/NzFDQE9R8fA3LA==" {
		// 验证成功，在当前这个用户的session data 保存一个键值对：isLogin=true
		// 1. 先从上下文中获取session data
		tmpSD, ok := ctx.Get(gin_session.SessionContextName)
		if !ok {
			panic("session middleware not support")
		}
		sd := tmpSD.(gin_session.SessionData)
		//2. 给session data设置isLogin = true
		sd.Set("isLogin", true)
		sd.Save() //调用Save，存储到数据库
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "login successfully",
			"uid":  111,
		})
		return
	} else { //验证失败，重新登陆
		ctx.JSON(http.StatusOK, gin.H{
			"code": -105,
			"msg":  "用户名或密码错误",
		})
		return
	}
}
func (c *LoginRouter) get_logout(ctx *gin.Context) {
	// 1. 先从上下文中获取session data
	tmpSD, _ := ctx.Get(gin_session.SessionContextName)
	sd := tmpSD.(gin_session.SessionData)
	//2. 给session data设置isLogin = true
	sd.Set("isLogin", false)
	sd.Save() //调用Save，存储到数据库
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "logout successfully",
	})
	return
}

func (c *LoginRouter) get_home(ctx *gin.Context) {
	uid, _ := ctx.Params.Get("uid")
	ctx.String(200, "get_home", uid)
}

func (c *LoginRouter) get_space(ctx *gin.Context) {
	sd, ok := AuthMiddleware(ctx)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "已登录",
			"uid":  111,
		})
		return
	} else {
		fmt.Printf("sd: %v\n", sd)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "未登录",
		})
	}
}

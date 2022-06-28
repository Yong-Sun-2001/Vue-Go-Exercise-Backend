package route

import (
	"bkd/middlewares/my_session"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRouter struct{}

func (c *LoginRouter) Router(engine *gin.Engine) {
	engine.POST("/login", c.post_login)
	engine.Use(my_session.AuthMiddleware()) //局部中间件
	{
		engine.GET("/space", c.get_space)
	}
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
	//验证输入的账号密码是否正确
	if now.Usr == "sun" && now.Pwd == "//I12vM/NzFDQE9R8fA3LA==" {
		// 1. 先从上下文中获取session data
		tmpSD, _ := ctx.Get(my_session.SessionContextName)
		sd := tmpSD.(my_session.SessionData)
		//2. 在session中设置isLogin，设置uid
		sd.Clear()
		sd.Set("isLogin", true)
		sd.Set("uid", 111)
		sd.Save() //调用Save，存储到数据库
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "login successfully",
			"uid":  111,
		})
	} else { //验证失败，重新登陆
		ctx.JSON(http.StatusOK, gin.H{
			"code": -105,
			"msg":  "用户名或密码错误",
		})
	}
}
func (c *LoginRouter) get_logout(ctx *gin.Context) {
	// 1. 先从上下文中获取session data
	tmpSD, _ := ctx.Get(my_session.SessionContextName)
	sd := tmpSD.(my_session.SessionData)
	//2. 清空会话
	sd.Clear()
	sd.Save() //调用Save，存储到数据库
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "logout successfully",
	})
}

//获取个人空间的个人信息
func (c *LoginRouter) get_space(ctx *gin.Context) {
	u := my_session.CurrentUser(ctx)
	if u != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "已登录",
			"uid":  111,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取用户信息错误",
		})
	}
}

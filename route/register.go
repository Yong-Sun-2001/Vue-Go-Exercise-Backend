package route

import (
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	// //基于my_session
	// a := LoginRouter{}
	// a.Router(e)
	//基于jwt
	a := LoginRouter_Jwt{}
	a.Router(e)
}

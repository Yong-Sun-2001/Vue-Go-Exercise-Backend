package route

import (
	"github.com/gin-gonic/gin"
)



func Register(e *gin.Engine) {
	a := LoginRouter{}
	a.Router(e)
}

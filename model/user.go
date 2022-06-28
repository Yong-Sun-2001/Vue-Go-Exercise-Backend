package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `json:"usr"`
	Password string `json:"pwd"`
	Nickname string
	Status   string
	Uid      uint64
	Avatar   string `gorm:"size:1000"`
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// GetUser 用ID获取用户
func GetUser(ID interface{}) (User, error) {
	// var user User
	// result := DB.First(&user, ID)
	// return user, result.Error
	return User{}, nil
}

// CurrentUser 获取当前用户
func CurrentUser(ctx *gin.Context) *User {
	if user, _ := ctx.Get("user"); user != nil {
		if u, ok := user.(*User); ok {
			fmt.Printf("u: %v\n", u)
			return u
		}
	}
	return nil
}

// SetPassword 设置密码
func (user *User) SetPassword(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}

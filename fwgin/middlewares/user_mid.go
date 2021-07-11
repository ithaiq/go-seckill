package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserMid struct {
}

func NewUserMid() *UserMid {
	return &UserMid{}
}

func (m *UserMid) OnRequest(context *gin.Context) error {
	fmt.Println("用户中间件")
	return nil
}

package classes

import (
	"github.com/gin-gonic/gin"
	"ithaiq/fwgin/core"
	"ithaiq/fwgin/models"
)

type UserClass struct {
}

func NewUserClass() *UserClass {
	return &UserClass{}
}

func (this *UserClass) GetUser(ctx *gin.Context) string {
	return "ok"
}

func (this *UserClass) GetDetail(ctx *gin.Context) core.IModel {
	return &models.UserModel{Id: 1, Name: "test"}
}

func (this *UserClass) GetList(ctx *gin.Context) core.SliceModel {
	list:=[]*models.UserModel{
		{Id: 1, Name: "test1"},
		{Id: 2, Name: "test2"},
	}
	return core.MakeSliceModel(list)
}

func (this *UserClass) Build(engine *core.Engine) {
	engine.Handle("GET", "/user", this.GetUser)
	engine.Handle("GET", "/user_detail", this.GetDetail)
	engine.Handle("GET", "/user_list", this.GetList)
}

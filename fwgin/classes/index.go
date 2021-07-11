package classes

import (
	"github.com/gin-gonic/gin"
	"ithaiq/fwgin/core"
)

type IndexClass struct {
}

func NewIndexClass() *IndexClass {
	return &IndexClass{}
}

func (this *IndexClass) GetIndex() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{"result": "Index success"})
	}
}

func (this *IndexClass) Build(engine *core.Engine) {
	engine.Handle("GET", "/index", this.GetIndex())
}

package core

import "github.com/gin-gonic/gin"

//Engine 代理gin.Engine
type Engine struct {
	*gin.Engine
	group *gin.RouterGroup
}

func NewEngine() *Engine {
	e := &Engine{Engine: gin.New()}
	e.Use(ErrorHandler())
	return e
}

func (e *Engine) Handle(httpMethod, relativePath string, handler interface{}) *Engine {
	//根据函数类型转换相应控制器
	if h := Convert(handler); h != nil {
		e.group.Handle(httpMethod, relativePath, h)
	}

	return e
}

func (e *Engine) Launch() {
	e.Run(":8080")
}
func (e *Engine) Mount(group string, classes ...IClass) *Engine {
	e.group = e.Group(group)
	for _, v := range classes {
		v.Build(e)
	}
	return e
}

func (e *Engine) Attach(m IMid) *Engine {
	e.Use(func(context *gin.Context) {
		err := m.OnRequest(context)
		if err != nil {
			context.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		} else {
			context.Next()
		}
	})
	return e
}

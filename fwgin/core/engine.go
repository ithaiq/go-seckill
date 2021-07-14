package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

//Engine 代理gin.Engine
type Engine struct {
	*gin.Engine
	group       *gin.RouterGroup
	beanFactory *BeanFactory
}

func NewEngine() *Engine {
	e := &Engine{Engine: gin.New(), beanFactory: NewBeanFactory()}
	e.Use(ErrorHandler())
	e.beanFactory.setBean(InitConfig()) // 配置文件加载进 bean 中
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
	var port int32 = 8080
	if config := e.beanFactory.GetBean(new(SysConfig)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	e.Run(fmt.Sprintf(":%d", port))
}
func (e *Engine) Mount(group string, classes ...IClass) *Engine {
	e.group = e.Group(group)
	for _, v := range classes {
		v.Build(e)
		e.beanFactory.inject(v)
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

func (e *Engine) Beans(beans ...interface{}) *Engine {
	e.beanFactory.setBean(beans...)
	return e

}

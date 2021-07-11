package core

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Responder interface {
	RespondTo() gin.HandlerFunc
}

var ResponderList []Responder

func init() {
	ResponderList = append(ResponderList,
		new(StringResponder),
		new(ModelResponder),
		new(SliceResponder))
}
func Convert(handler interface{}) gin.HandlerFunc {
	hRef := reflect.ValueOf(handler)
	for _, r := range ResponderList {
		rRef := reflect.ValueOf(r).Elem()
		if hRef.Type().ConvertibleTo(rRef.Type()) {
			rRef.Set(hRef)
			return rRef.Interface().(Responder).RespondTo()
		}
	}
	return nil
}

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, this(context))
	}
}

type ModelResponder func(*gin.Context) IModel

func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

type SliceResponder func(*gin.Context) SliceModel

func (this SliceResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(this(context)))
	}
}

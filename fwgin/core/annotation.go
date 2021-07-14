package core

import (
	"fmt"
	"reflect"
	"strings"
)

type Annotation interface {
	SetTag(tag reflect.StructTag)
	String() string
}

var AnnotationList []Annotation

func IsAnnotation(t reflect.Type) bool {
	for _, item := range AnnotationList {
		if reflect.TypeOf(item) == t {
			return true
		}
	}
	return false
}


// init 包构造函数
func init() {
	AnnotationList = make([]Annotation, 0)
	AnnotationList = append(AnnotationList, new(Value))
}


type Value struct {
	tag reflect.StructTag
	BeanFactory *BeanFactory
}

func (v *Value) SetTag(tag reflect.StructTag) {
	v.tag = tag
}

func (this *Value) String() string {
	get_prefix := this.tag.Get("prefix")
	if get_prefix == "" {
		return ""
	}
	prefix := strings.Split(get_prefix, ".")
	if config := this.BeanFactory.GetBean(new(SysConfig)); config != nil {
		get_value := GetConfigValue(config.(*SysConfig).Config, prefix, 0)
		if get_value != nil {
			return fmt.Sprintf("%v", get_value)
		} else {
			return ""
		}
	} else {
		return ""
	}
}
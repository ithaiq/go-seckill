package core

import "reflect"

type BeanFactory struct {
	beans []interface{}
}

func NewBeanFactory() *BeanFactory {
	bf := &BeanFactory{beans: make([]interface{}, 0)}
	bf.beans = append(bf.beans, bf)
	return bf
}

// GetBean 外部使用获取注入的属性
func (e *BeanFactory) GetBean(bean interface{}) interface{} {
	return e.getBean(reflect.TypeOf(bean))
}

func (e *BeanFactory) getBean(t reflect.Type) interface{} {
	for _, p := range e.beans {
		if t == reflect.TypeOf(p) {
			return p
		}
	}
	return nil
}

func (e *BeanFactory) setBean(beans ...interface{}) {
	e.beans = append(e.beans, beans...)
}

// Inject 给外部用的 （后面还要改,这个方法不处理注解)
func (e *BeanFactory) Inject(object interface{}) {
	vObject := reflect.ValueOf(object)
	if vObject.Kind() == reflect.Ptr { //由于不是控制器 ，所以传过来的值 不一定是指针。因此要做判断
		vObject = vObject.Elem()
	}
	for i := 0; i < vObject.NumField(); i++ {
		f := vObject.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if p := e.getBean(f.Type()); p != nil && f.CanInterface() {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
		}
	}
}

// inject 把bean注入到控制器中 (内部方法,用户控制器注入。并同时处理注解)
func (e *BeanFactory) inject(class IClass) {
	vClass := reflect.ValueOf(class).Elem()
	vClassT := reflect.TypeOf(class).Elem()
	for i := 0; i < vClass.NumField(); i++ {
		f := vClass.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if IsAnnotation(f.Type()) {
			f.Set(reflect.New(f.Type().Elem()))
			f.Interface().(Annotation).SetTag(vClassT.Field(i).Tag)
			e.Inject(f.Interface())
			continue
		}
		if p := e.getBean(f.Type()); p != nil {
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())
		}
	}
}
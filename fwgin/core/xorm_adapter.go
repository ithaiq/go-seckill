package core

import (
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type XormAdapter struct {
	*xorm.Engine
}

func NewXormAdapter() *XormAdapter {
	return &XormAdapter{}
}

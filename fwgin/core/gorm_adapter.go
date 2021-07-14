package core

import (
	"fmt"
	"gorm.io/gorm"
)

type GormAdapter struct {
	*gorm.DB
}

func NewGormAdapter() *GormAdapter {
	return &GormAdapter{}
}

func (g *GormAdapter) Test(str string) {
	fmt.Println(str)
}
package main

import (
	"ithaiq/fwgin/classes"
	"ithaiq/fwgin/core"
	"ithaiq/fwgin/middlewares"
)

func main() {
	core.NewEngine().
		Beans(core.NewGormAdapter()).
		Attach(middlewares.NewUserMid()).
		Mount("v1", classes.NewIndexClass()).
		Mount("v2", classes.NewUserClass()).
		Launch()
}
package core

import (
	"encoding/json"
	"log"
)

type IModel interface {
	String() string
}

type SliceModel string

func MakeSliceModel(model interface{}) SliceModel {
	str, err := json.Marshal(model)
	if err != nil {
		log.Println(err)
	}
	return SliceModel(str)
}

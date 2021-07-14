package core

import (
	"io/ioutil"
	"os"
)

func LoadConfigFile() ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file := dir + "/application.yaml"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}

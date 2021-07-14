package core

import (
	"gopkg.in/yaml.v2"
	"log"
)

type ServerConfig struct {
	Port int32
	Name string
}

// SysConfig 系统配置
type SysConfig struct {
	Server *ServerConfig
	Config UserConfig
}

type UserConfig map[interface{}]interface{}

// NewSysConfig 初始化默认配置
func NewSysConfig() *SysConfig {
	return &SysConfig{Server: &ServerConfig{Port: 8080, Name: "fwgin"}}
}

func InitConfig() *SysConfig {
	config := NewSysConfig()
	content, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func GetConfigValue(m UserConfig, prefix []string, index int) interface{} {
	key := prefix[index]
	if v, ok := m[key]; ok {
		if index == len(prefix)-1 {
			return v
		} else {
			index = index + 1
			if mv, ok := v.(UserConfig); ok {
				return GetConfigValue(mv, prefix, index)
			} else {
				return nil
			}
		}
	}
	return nil
}

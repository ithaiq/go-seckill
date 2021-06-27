package utils

import (
	"encoding/json"
	"io/ioutil"
	"ithaiq/gtcp/giface"
)

type GlobalObj struct {
	TcpServer      giface.IServer
	Host           string
	TcpPort        int
	Name           string
	Version        string
	MaxConn        int
	MaxPackageSize uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8888,
		Name:           "GtcpServerApp",
		Version:        "v3.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/gtcp.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
	GlobalObject.Reload()
}

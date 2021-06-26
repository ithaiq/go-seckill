package gnet

import (
	"ithaiq/gtcp/giface"
	"testing"
)

var _ giface.IRouter = (*PingRouter)(nil)

type PingRouter struct {
	BaseRouter
}

func (p *PingRouter) Handle(request giface.IRequest) {
	request.GetConnection().Send(request.GetData())
}

func (p *PingRouter) PostHandle(request giface.IRequest) {
	request.GetConnection().Send(request.GetData())
}

func TestNewServer(t *testing.T) {
	server := NewServer("gtcp")
	server.AddRouter(&PingRouter{})
	server.Serve()
}

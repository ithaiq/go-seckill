package gnet

import "ithaiq/gtcp/giface"

var _ giface.IRouter = (*BaseRouter)(nil)
type BaseRouter struct {

}

func (b *BaseRouter) PreHandle(request giface.IRequest) {
}

func (b *BaseRouter) Handle(request giface.IRequest) {
}

func (b *BaseRouter) PostHandle(request giface.IRequest) {
}


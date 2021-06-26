package gnet

import "ithaiq/gtcp/giface"

var _ giface.IRequest = (*Request)(nil)

type Request struct {
	conn giface.IConnection
	data []byte
}

func (r *Request) GetConnection() giface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}

package gnet

import (
	"ithaiq/gtcp/giface"
	"net"
)

var _ giface.IConnection = (*Connection)(nil)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	//handleAPI giface.HandleFunc
	ExitChan chan bool
	Router   giface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router giface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) Start() {
	go c.StartReader()
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	_, err := c.Conn.Write(data)
	return err
}

func (c *Connection) StartReader() {
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			continue
		}
		/*if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			break
		}*/
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request giface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

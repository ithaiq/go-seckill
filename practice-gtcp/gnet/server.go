package gnet

import (
	"fmt"
	"ithaiq/gtcp/giface"
	"net"
)

var _ giface.IServer = (*Server)(nil)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	if _, err := conn.Write(data[:cnt]); err != nil {
		return err
	}
	return nil
}

func (s *Server) Start() {
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			return
		}
		var cid uint32
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				continue
			}
			newConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			go newConn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func NewServer(name string) giface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8888,
	}
	return s
}

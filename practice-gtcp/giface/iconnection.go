package giface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	Send(data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error

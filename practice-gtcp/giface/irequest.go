package giface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}
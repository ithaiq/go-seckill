package giface

type IServer interface {
	Start()
	Stop()
	Serve()
}

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ithaiq/grpc/codec"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 1024

/*服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType 解码剩余的内容*/
//| Option | Header1 | Body1 | Header2 | Body2 | ...

type Option struct {
	MagicNumber    int
	CodecType      codec.Type
	ConnectTimeout time.Duration //连接超时
	HandleTimeout  time.Duration //处理超时
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func (server *Server) ServeConn(conn net.Conn) {
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("json decode options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn), &opt)
}

var invalidRequest = struct{}{}

//serveCodec 开始正式处理请求
func (server *Server) serveCodec(cc codec.Codec, opt *Option) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		//处理请求是并发的，但是回复请求的报文必须是逐个发送的(使用锁控制)
		go server.handleRequest(cc, req, sending, wg, opt.HandleTimeout)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) Register(rcvR interface{}) error {
	s := newService(rcvR)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("service already defined: " + s.name)
	}
	return nil
}

func (server *Server) findService(serviceMethod string) (svc *service, mType *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svcL, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("can't find service " + serviceName)
		return
	}
	svc = svcL.(*service)
	mType = svc.method[methodName]
	if mType == nil {
		err = errors.New("can't find method " + methodName)
	}
	return
}

func Register(rcvR interface{}) error {
	return DefaultServer.Register(rcvR)
}

type request struct {
	h            *codec.Header
	argV, replyV reflect.Value
	mType        *methodType
	svc          *service
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.svc, req.mType, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argV = req.mType.newArgv()
	req.replyV = req.mType.newReplyV()

	argvi := req.argV.Interface()
	if req.argV.Type().Kind() != reflect.Ptr {
		argvi = req.argV.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("read body err:", err)
		return req, err
	}
	return req, nil
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	called := make(chan struct{})
	sent := make(chan struct{})
	go func() {
		err := req.svc.call(req.mType, req.argV, req.replyV)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyV.Interface(), sending)
		sent <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-sent
		return
	}
	select {
	case <-time.After(timeout): //time.After() 先于 called 接收到消息，说明处理已经超时
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called: //called 信道接收到消息，代表处理没有超时，继续执行 sendResponse
		<-sent
	}
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("write response error:", err)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"ithaiq/grpc/codec"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 1024

/*服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType 解码剩余的内容*/
//| Option | Header1 | Body1 | Header2 | Body2 | ...

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

type Server struct{}

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
	server.serveCodec(f(conn))
}

var invalidRequest = struct{}{}

//serveCodec 开始正式处理请求
func (server *Server) serveCodec(cc codec.Codec) {
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
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

type request struct {
	h            *codec.Header
	argV, replyV reflect.Value
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
	req.argV = reflect.New(reflect.TypeOf(""))
	//读取请求体
	if err = cc.ReadBody(req.argV.Interface()); err != nil {
		log.Println("read argV err:", err)
		return req, err
	}
	return req, nil
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("接收到：",req.h, req.argV.Elem())
	//设置响应内容
	req.replyV = reflect.ValueOf(fmt.Sprintf("grpc rsp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyV.Interface(), sending)
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("write response error:", err)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

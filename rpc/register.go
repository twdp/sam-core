package rpc

import "github.com/smallnest/rpcx/server"

type RpcMeta struct {
	Name string
	Rcvr interface{}
	Meta string
}

var rpcContainer []*RpcMeta

func Register(name string, rcvr interface{}, metadata string) {
	m := &RpcMeta{
		Name: name,
		Rcvr: rcvr,
		Meta: metadata,
	}
	rpcContainer = append(rpcContainer, m)
}


func RegisterToRpcx(s *server.Server) {
	for _, rpc := range rpcContainer {
		s.RegisterName(rpc.Name, rpc.Rcvr, rpc.Meta)
	}
}
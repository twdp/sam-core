// +build !zookeeper

package rpc

import (
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
)

func AddRegistryPlugin(s *server.Server) {
	r := client.InprocessClient
	s.Plugins.Add(r)
}
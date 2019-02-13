// +build zookeeper

package rpc

import (
	"github.com/hzwy23/hauth/utils/logs"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"tianwei.pro/sam-core/config"
	"time"
)

func AddRegistryPlugin(s *server.Server) {
	r := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@" + config.Conf.DefaultString("addr", "localhost:29998"),
		ZooKeeperServers: config.Conf.DefaultStrings("zkAddr", []string{ "localhost:2181" }),
		BasePath:         config.Conf.DefaultString("basePath", "sam"),
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   time.Minute,
	}
	err := r.Start()
	if err != nil {
		logs.Error(err)
	}
	s.Plugins.Add(r)
}
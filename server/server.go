package main

import (
	"flag"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/server"
	"tianwei.pro/sam-core/config"
	_ "tianwei.pro/sam-core/config"
	_ "tianwei.pro/sam-core/impl"
	_ "tianwei.pro/sam-core/model"
	"tianwei.pro/sam-core/rpc"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", config.Conf.String("db.url"), 30)

	// create table
	orm.RunSyncdb("default", false, true)

	if config.Conf.DefaultString("runmode", "prod") != "prod" {
		orm.Debug = true
	}
}

func main() {
	flag.Parse()
	s := server.NewServer()
	rpc.AddRegistryPlugin(s)

	rpc.RegisterToRpcx(s)
	s.Serve("tcp", config.Conf.DefaultString("addr", "localhost:29998"))
}

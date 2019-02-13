// +build zookeeper

package config

import (
	"flag"
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config/yaml"
	"github.com/astaxie/beego/logs"
)

var (
	ConfigName = flag.String("c", "conf/conf.yaml", "配置文件地址")
)

func init() {
	c, err := config.NewConfig("yaml", *ConfigName)
	if err != nil {
		logs.Error("parse config failed. %v", err)
		panic(err)
	}
	Conf = c
}

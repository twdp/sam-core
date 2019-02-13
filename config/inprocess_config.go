// +build !zookeeper

package config

import "github.com/astaxie/beego"

func init() {
	Conf = beego.AppConfig
}
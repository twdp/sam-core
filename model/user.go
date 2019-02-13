package model

import "tianwei.pro/business/model"

type User struct {

	model.Base

	UserName string `orm:"size(64);unique"`

	DisplayName string `orm:"size(64)"`

	Avatar string

	Email string `orm:"size(64);unique"`

	Phone string `orm:"size(64);unique"`

	Sex int8

	Password string

	Type int8

	Status int8

	// 退出时记录哪些端需要重新登录
	NeedLoginTerminus string `orm:"size(120)"`
}

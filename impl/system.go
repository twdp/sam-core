package impl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"tianwei.pro/business"
	"tianwei.pro/sam-core/dto/req"
	"tianwei.pro/sam-core/dto/res"
	"tianwei.pro/sam-core/model"
	"tianwei.pro/sam-core/repository"
	"tianwei.pro/sam-core/rpc"
	"time"
)

var (
	DuplicationKeyErr = errors.New("名称重复")
)

func init() {
	f := &SystemFacadeImpl{
		systemRepository: repository.SystemRepositoryInstance,
	}
	rpc.Register("SystemFacadeImpl", f, "")
}

type SystemFacadeImpl struct {
	//sync.Mutex

	systemRepository *repository.SystemRepository
}

func (s *SystemFacadeImpl) Stay(context context.Context, stay *req.SystemStay, reply *res.StayResponse) error {
	fmt.Println(stay)
	system := &model.System{
		Name:        stay.Name,
		Avatar:      stay.Avatar,
		Description: stay.Description,
		Strategy:    stay.Strategy,
		AppKey:      s.generateAppKey(),
		Secret:      s.generateAppKey(),
	}

	if err := s.systemRepository.Stay(system, stay.Operator); err != nil {
		if business.IsDuplicationKeyError(err) {
			logs.Warn("stay system duplication. stay: %v, err: %v", stay, err)
			return reply.Error(DuplicationKeyErr)
		} else {
			logs.Error("stay system error: stay: %v, err: %v", stay, err)
			return reply.Error(SystemError)
		}
	}

	// 创建系统管理员
	// 关联管理员

	reply.Name = system.Name
	reply.Description = system.Description
	reply.Id = system.Id
	reply.Avatar = system.Avatar
	reply.AppKey = system.AppKey
	reply.Secret = system.Secret
	reply.Strategy = system.Strategy
	reply.Success = true

	return nil
}

func (s *SystemFacadeImpl) ListByOwner(context context.Context, owner int64, reply *res.SystemListResponse) error {

	var uroles []*model.UserRole
	orm.NewOrm().QueryTable(&model.UserRole{}).Filter("UserId", owner).Filter("RoleId", 1).All(&uroles)

	if len(uroles) == 0 {
		return nil
	}
	var ss []*model.System
	var idss []int64
	for _, s := range uroles {
		idss = append(idss, s.SystemId)
	}
	orm.NewOrm().QueryTable(&model.System{}).Filter("id__in", idss).All(&ss)

	var list []*res.SystemList
	for _, s := range ss {
		list = append(list, &res.SystemList{
			Id:                  s.Id,
			CreatedAt:           s.CreatedAt.Value(),
			UpdatedAt:           s.UpdatedAt.Value(),
			Status:              s.Status,
			Name:                s.Name,
			Description:         s.Description,
			Avatar:              s.Avatar,
			AppKey:              s.AppKey,
			Secret:              s.Secret,
			KeepSign:            s.KeepSign,
			Strategy:            s.Strategy,
			TemplateRoleVisible: s.TemplateRoleVisible,
		})
	}
	reply.List = list

	return nil
}

func (s *SystemFacadeImpl) generateAppKey() (result string) {
	//s.Lock()
	//defer s.Unlock()

	result = business.CastInt64ToString(time.Now().Unix())

	randLength := 20
	randType := "Aa0"
	var num = "0123456789"
	var lower = "abcdefghijklmnopqrstuvwxyz"
	var upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := bytes.Buffer{}
	if strings.Contains(randType, "0") {
		b.WriteString(num)
	}
	if strings.Contains(randType, "a") {
		b.WriteString(lower)
	}
	if strings.Contains(randType, "A") {
		b.WriteString(upper)
	}
	var str = b.String()
	var strLen = len(str)
	if strLen == 0 {
		result += ""
		return
	}

	rand.Seed(time.Now().UnixNano())
	b = bytes.Buffer{}
	for i := 0; i < randLength; i++ {
		b.WriteByte(str[rand.Intn(strLen)])
	}
	result += b.String()
	return
}

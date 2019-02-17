package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"tianwei.pro/business"
	"tianwei.pro/sam-agent"
	"tianwei.pro/sam-core/model"
	"tianwei.pro/sam-core/repository"
	"tianwei.pro/sam-core/rpc"
)

var (
	agent               *SamCoreAgentImpl
	AppKeyOrSecretError = errors.New("请检查appKey或secret")
	SystemError         = errors.New("权限系统错误")
)

type SamCoreAgentImpl struct {
	systemRepository *repository.SystemRepository

	apiRepository *repository.ApiRepository
}

func init() {
	agent = &SamCoreAgentImpl{
		systemRepository: repository.SystemRepositoryInstance,
		apiRepository:    repository.ApiRepositoryInstance,
	}
	rpc.Register("SamAgentFacadeImpl", agent, "")
}

func (s *SamCoreAgentImpl) LoadSystemInfo(ctx context.Context, param *sam_agent.SystemInfoParam, reply *sam_agent.SystemInfo) error {
	if systemInfo, err := s.verifySecret(param.AppKey, param.Secret); err != nil {
		return reply.Error(err)
	} else {
		routes, err := s.apiRepository.FindUrlBySystemId(systemInfo.Id)
		if business.IsError(err) {
			return reply.Error(err)
		}
		var apis []*sam_agent.Router

		for _, api := range routes {
			apis = append(apis, &sam_agent.Router{
				Id:     api.Id,
				Url:    api.Path,
				Method: api.Method,
				Type:   api.VerificationType,
			})
		}

		reply.Id = systemInfo.Id
		reply.PermissionType = systemInfo.Strategy
		reply.KeepSign = systemInfo.KeepSign
		reply.Routers = apis

		return nil
	}
}

func (s *SamCoreAgentImpl) verifySecret(appKey, secret string) (*model.System, error) {
	if system, err := s.systemRepository.FindByAppKey(appKey); err != nil {
		logs.Warn("app key: %v not found", appKey)
		return nil, AppKeyOrSecretError
	} else {
		if system.Secret != secret {
			logs.Warn("param: %v, system info: %v", secret, system)
			return nil, AppKeyOrSecretError
		} else {

		}
		return system, nil
	}
}

func (s *SamCoreAgentImpl) VerifyToken(ctx context.Context, param *sam_agent.VerifyTokenParam, reply *sam_agent.UserInfo) error {
	system, err := s.verifySecret(param.AppKey, param.Secret)
	if err != nil {
		return reply.Error(err)
	}
	if user, err := tokenFacadeImpl.DecodeToken(param.Token); err != nil {
		return reply.Error(err)
	} else if err := s.loadUserInfo(user.Id, reply); err != nil {
		return reply.Error(err)
	}

	// todo:: add cache

	var userRoles []*model.UserRole
	if _, err := orm.NewOrm().QueryTable(&model.UserRole{}).Filter("SystemId", system.Id).Filter("UserId", reply.Id).All(&userRoles); err != nil {
		logs.Error("read user role mapping failed. userInfo: %v, system: %v, err: %v", reply, system, err)
		return reply.Error(SystemError)
	}
	if len(userRoles) == 0 {
		return nil
	}

	fmt.Println(system)
	return nil
}

func (s *SamCoreAgentImpl) loadUserInfo(userId int64, reply *sam_agent.UserInfo) error {
	u := &model.User{}
	u.Id = userId
	if err := orm.NewOrm().Read(u); err != nil {
		logs.Error("find user by token failed. err: %v", err)
		return reply.Error(SystemError)
	}
	reply.Id = u.Id
	reply.UserName = u.UserName
	reply.Avatar = u.Avatar
	reply.Email = u.Email
	reply.Phone = u.Phone

	return nil
}

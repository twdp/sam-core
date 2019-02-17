package impl

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	"tianwei.pro/sam-agent"
	"tianwei.pro/sam-core/const"
	"tianwei.pro/sam-core/dto"
	"tianwei.pro/sam-core/dto/req"
	"tianwei.pro/sam-core/dto/res"
	"tianwei.pro/sam-core/repository"
	"tianwei.pro/sam-core/rpc"
)

var (
	EmailOrPassErr = errors.New("账号或密码错误")
	UserNotActive  = errors.New("用户未激活或被冻结")
)

func init() {
	f := &UserFacadeImpl{
		userRepository: repository.UserRepositoryInstance,
	}
	rpc.Register("UserFacadeImpl", f, "")
}

type UserFacadeImpl struct {
	userRepository *repository.UserRepository
}

func (u *UserFacadeImpl) Login(ctx context.Context, loginParam *req.EmailLoginDto, reply *res.LoginDto) error {
	if user, err := u.userRepository.FindByEmail(loginParam.Email); err != nil {
		return reply.Error(err)
	} else if _, err := business.ValidateCrypto(loginParam.Password, user.Password); err != nil {
		return EmailOrPassErr
	} else {
		if user.Status != _const.Active {
			logs.Warn("user status not active, user: %s, orm user: %v", loginParam.Email, user)
			return reply.Error(UserNotActive)
		}
		userDto := &dto.UserDto{
			BaseDto: dto.BaseDto{
				Id: user.Id,
			},
			UserName: user.UserName,
		}
		if token, err := tokenFacadeImpl.EncodeToken(userDto, loginParam.Terminal); err != nil {
			return reply.Error(err)
		} else {
			reply.Token = token

			r := &sam_agent.UserInfo{}
			if err := agent.VerifyToken(context.Background(), &sam_agent.VerifyTokenParam{
				SystemInfoParam: sam_agent.SystemInfoParam{
					AppKey: loginParam.AppKey,
					Secret: loginParam.Secret,
				},
				Token: token,
			}, r); err != nil {
				logs.Error("verify token failed. loginParam: %v, err: %v", loginParam, err)
				return SystemError
			}
			reply.UserInfo = *r

			return nil
		}
	}
}

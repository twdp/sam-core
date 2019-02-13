package impl

import (
	"context"
	"github.com/pkg/errors"
	"tianwei.pro/business"
	"tianwei.pro/sam-core/dto"
	"tianwei.pro/sam-core/dto/req"
	"tianwei.pro/sam-core/dto/res"
	"tianwei.pro/sam-core/repository"
	"tianwei.pro/sam-core/rpc"
)

var (
	EmailOrPassErr = errors.New("账号或密码错误")
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
		return err
	} else if _, err := business.ValidateCrypto(loginParam.Password, user.Password); err != nil {
		return EmailOrPassErr
	} else {
		userDto := &dto.UserDto{
			BaseDto: dto.BaseDto{
				Id: user.Id,
			},
			UserName: user.UserName,
		}
		if token, err := tokenFacadeImpl.EncodeToken(userDto, loginParam.Terminal); err != nil {
			return err
		} else {
			reply.Token = token
			// todo:  sam_agent user info
			return nil
		}
	}
}

package service

import (
	"context"
	"errors"
	"github/yyfzy/mybook/internal/domain"
	"github/yyfzy/mybook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱不正确")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	//存储
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// debug
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) UpdateBasicInfo(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 快路径，先查询，如果存在则直接返回
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// nil 与 不为ErrUserNotFound 均会进入此分支
		return u, err
	}
	// 在系统资源不足，触发降级之后不执行慢路径
	if ctx.Value("降级") == "true" {
		return domain.User{}, errors.New("系统降级，禁止创建用户")
	}
	// 慢路径
	u = domain.User{Phone: phone}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里存在主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}

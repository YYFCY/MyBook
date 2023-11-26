package repository

import (
	"context"
	"fmt"
	"github/yyfzy/mybook/internal/domain"
	"github/yyfzy/mybook/internal/repository/cache"
	"github/yyfzy/mybook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{dao: dao, cache: c}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: u.Id, Email: u.Email, Password: u.Password}, nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	err := r.dao.Update(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	})
	if err != nil {
		return err
	}
	return r.cache.Delete(ctx, u.Id)
}

// FindById 如果没有数据，返回一个特定的error； error为nil，就认为缓存里有数据
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从cache找，再从dao里找，找到了回写cache
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		//
		return u, nil
	}
	//if err == cache.ErrKeyNotExit {
	//	// 去数据库加载
	//}
	// 如果err 为其他错误，比如 err = io.EOF ,要不要去数据库加载？
	// 选不加载，用户体验差；选择加载--做好兜底，万一redis宕机，要保护住数据库，可以用数据库限流，用orm的middleware （面试可用）；

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
		Nickname: ue.Nickname,
		Birthday: ue.Birthday,
		AboutMe:  ue.AboutMe,
		Ctime:    time.UnixMilli(ue.Ctime),
	}
	err = r.cache.Set(ctx, u)
	if err != nil {
		fmt.Println(err)

	}
	return u, nil
}

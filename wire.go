//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github/yyfzy/mybook/internal/repository"
	"github/yyfzy/mybook/internal/repository/cache"
	"github/yyfzy/mybook/internal/repository/dao"
	"github/yyfzy/mybook/internal/service"
	"github/yyfzy/mybook/internal/web"
	"github/yyfzy/mybook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

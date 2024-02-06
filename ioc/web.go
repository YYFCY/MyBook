package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github/yyfzy/mybook/internal/web"
	"github/yyfzy/mybook/internal/web/middleware"
	"github/yyfzy/mybook/pkg/ginx/ratelimit"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		corsHandle(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/login", "/users/signup", "/users/login_sms/code/send", "/users/login_sms").
			Build(),
	}
}

func corsHandle() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    nil,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "abc.com")
		},
		AllowMethods:           nil,
		AllowHeaders:           nil,
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 0,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	})
}

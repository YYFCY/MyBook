package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	ginRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github/yyfzy/mybook/config"
	"github/yyfzy/mybook/internal/repository"
	"github/yyfzy/mybook/internal/repository/cache"
	"github/yyfzy/mybook/internal/repository/dao"
	"github/yyfzy/mybook/internal/service"
	"github/yyfzy/mybook/internal/service/sms/memory"
	"github/yyfzy/mybook/internal/web"
	"github/yyfzy/mybook/internal/web/middleware"
	"github/yyfzy/mybook/pkg/ginx/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisClient := initRedis()
	server := initWebServer(redisClient)
	u := initUser(db, redisClient)
	u.RegisterRoutes(server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, yyf")
	})
	server.Run(":8080")

}

func initWebServer(client redis.Cmdable) *gin.Engine {
	server := gin.Default()
	server.Use(ratelimit.NewBuilder(client, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
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
	}))

	//store := cookie.NewStore([]byte("secret"))
	store, err := ginRedis.NewStore(16, "tcp", config.Config.Redis.Addr, "",
		[]byte("3o4q6EshoibpRdTB6iPCayquqFmMQzkv"), []byte("naspBhPdXGTMOG9OoRaIukf48sf8WUXU"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login", "/users/signup", "/users/login_sms/code/send", "/users/login_sms").
		Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:13309)/webook"))

	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	// 这里演示读取特定的某个字段
	cmd := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return cmd
}

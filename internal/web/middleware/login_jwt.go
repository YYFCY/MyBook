package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"mybook/internal/web"
	"net/http"
	"strings"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(paths ...string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, paths...)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		claims := &web.UserClaims{}
		tokenHeaderStr := ctx.GetHeader("Authorization")
		if tokenHeaderStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segments := strings.SplitN(tokenHeaderStr, " ", 2)
		tokenStr := segments[1]
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("3o4q6EshoibpRdTB6iPCayquqFmMQzkv"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println(token.Valid)
		if claims == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 存在安全问题, 需要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)
	}
}

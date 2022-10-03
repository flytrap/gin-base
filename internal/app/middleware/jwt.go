package middleware

import (
	"errors"
	"net/http"

	"github.com/flytrap/gin-base/internal/app/services"
	"github.com/flytrap/gin-base/pkg/util"
	"github.com/gin-gonic/gin"
)

func JWTAuth(as services.AuthService, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := as.GetToken(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, util.ErrorWarper(errors.New("unauthorized")))
			c.Abort()
			return
		}

		userId, err := as.ParseUserID(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, util.ErrorWarper(err))
			c.Abort()
			return
		}
		c.Set("userID", userId)
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/flytrap/gin-base/pkg/logger"
	"github.com/gin-gonic/gin"
)

type SkipperFunc func(*gin.Context) bool

func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

// errorHttp 统一500错误处理函数
func ErrorHttp(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			logger.Printf("panic: %v\n", r)
			debug.PrintStack()

			c.JSON(http.StatusOK, gin.H{
				"title": "500",
			})
		}
	}()
	c.Next()
}

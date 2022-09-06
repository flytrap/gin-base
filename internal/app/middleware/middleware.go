package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// errorHttp 统一500错误处理函数
func ErrorHttp(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			logger.Printf("panic: %v\n", r)
			debug.PrintStack()

			c.HTML(http.StatusOK, "500.html", gin.H{
				"title": "500",
			})
		}
	}()
	c.Next()
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var _ IRouter = (*Router)(nil)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
} // end

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	a.RegisterPage(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
}

// RegisterPage register page group router
func (a *Router) RegisterPage(app *gin.Engine) {

	app.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.Redirect(http.StatusSeeOther, "/")
	})
}

package app

import (
	"github.com/flytrap/gin_template/docs"
	"github.com/flytrap/gin_template/internal/app/config"
	"github.com/flytrap/gin_template/internal/app/middleware"
	"github.com/flytrap/gin_template/internal/app/router"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	docs.SwaggerInfo.BasePath = "/api/v1"

	app.LoadHTMLGlob("templates/**/*")
	app.Static("/assets", "assets")

	// Swagger
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	app.Use(middleware.ErrorHttp)
	// Router register
	r.Register(app)

	return app
}

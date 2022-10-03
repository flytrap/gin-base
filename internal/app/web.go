package app

import (
	"github.com/flytrap/gin-base/docs"
	"github.com/flytrap/gin-base/internal/app/config"
	"github.com/flytrap/gin-base/internal/app/middleware"
	"github.com/flytrap/gin-base/internal/app/router"
	"github.com/flytrap/gin-base/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerValidation() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validate.RegisterValidation("birth", util.CheckBirthDate)
		validate.RegisterValidation("phone", util.MobileValidator)
		validate.RegisterValidation("idCard", util.IdCardValidator)
	}
}

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.Default()
	registerValidation()

	prefixes := r.Prefixes()

	// logger
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))
	// CORS
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	docs.SwaggerInfo.BasePath = "/api/v1"

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

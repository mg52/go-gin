package app

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-gin/app/api"
	"go-gin/config"
	"go-gin/db"
	"go-gin/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "go-gin/docs"
)

type Routes struct {
}

func (app Routes) StartGin() {
	config.SetConfigFile("config", "./config", "json")

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	//r.GET("swagger/*any", middlewares.NewSwagger())
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	publicRoute := r.Group("/api/v1")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	r.Static("/template/css", "./template/css")
	r.Static("/template/images", "./template/images")
	//r.Static("/template", "./template")

	r.NoRoute(func(context *gin.Context) {
		//context.File("./template/route_not_found.html")
		context.File("./template/index.html")
	})

	api.ApplyToDoAPI(publicRoute, resource)
	api.ApplyUserAPI(publicRoute, resource)
	r.Run(":" + config.GetString("host.port"))
}

package app

import (
	"mg52-gin/app/api"
	"mg52-gin/config"
	"mg52-gin/db"
	"mg52-gin/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Routes struct {
}

func (app Routes) StartGin() {
	config.SetConfigFile("config", "./config", "json")

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", middlewares.NewSwagger())

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

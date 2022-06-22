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
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	publicRoute := r.Group("/api/v1")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	api.ApplyToDoAPI(publicRoute, resource)
	api.ApplyUserAPI(publicRoute, resource)
	r.Run(":" + config.GetString("host.port"))
}

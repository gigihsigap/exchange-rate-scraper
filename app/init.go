package app

import (
	"backend/app/api"
	"backend/db"
	"backend/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Routes struct {
}

func (app Routes) StartGin() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	// r.GET("swagger/*any",middlewares.NewSwagger())

	publicRoute := r.Group("/api")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	// r.Static("/template/css", "./template/css")
	// r.Static("/template/images", "./template/images")
	//r.Static("/template", "./template")

	r.NoRoute(func(context *gin.Context) {
		//context.File("./template/route_not_found.html")
		//context.File("./template/index.html")
	})

	api.ApplyEntryAPI(publicRoute, resource)
	r.Run(":" + os.Getenv("PORT"))
}

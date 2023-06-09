package routes

import (
	"apiingolang/activity/business/routes"
	"apiingolang/activity/db"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	api := router.Group("/api")
	publicGroup := api.Group("/public")

	routes.ActivityRoutes(publicGroup, db.ClientActivity)
}

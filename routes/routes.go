package routes

import (
	"apiingolang/activity/business/routes"
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/db"
	"apiingolang/activity/middleware"
	"apiingolang/activity/middleware/corel"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.Use(corel.DefaultGinHandlers...)
	// adding recovery for api flow
	api.Use(middleware.Recovery(logging.Logger))
	api.Use(logging.Logger.Gin())
	publicGroup := api.Group("/public")

	routes.ActivityRoutes(publicGroup, db.ClientActivity)
}

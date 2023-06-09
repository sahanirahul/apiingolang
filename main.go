package main

import (
	"apiingolang/activity/routes"
	"net/http"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) { c.JSON(http.StatusOK, "OK") }

func main() {
	router := gin.Default()
	router.GET("/health", health)

	routes.InitRoutes(router)

	err := router.Run(":" + os.Getenv("PORT"))

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

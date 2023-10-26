package main

import (
	"conectivity-checker-wizard/controllers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// create a router
	router := gin.Default()

	// create and set up a session middleware
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// load all the templates
	router.LoadHTMLGlob("templates/*")

	// Serve static files (e.g., css, js, images)
	router.Static("/static/", "./static/")

	// configure routes
	mainController := new(controllers.MainController)
	router.GET("/", mainController.Home)
	router.POST("/api/rule/:ruleID", mainController.Execute)
	router.GET("/cilium", mainController.CiliumPolicies)

	router.Run(":8080")
}

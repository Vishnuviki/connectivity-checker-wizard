package main

import (
	"os"

	"conectivity-checker-wizard/config"
	"conectivity-checker-wizard/controllers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// create a router
	router := gin.New()
	gin.SetMode(os.Getenv("GIN_MODE"))
	// gin.SetMode(gin.ReleaseMode) // Modes - "debug", "release"

	// create and set up a session middleware
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// load all the templates
	router.LoadHTMLGlob("views/templates/*")

	// serve static files (e.g., css, js, images)
	router.Static("/views/static/", "./views/static/")

	// configure routes
	mainController := new(controllers.MainController)
	router.GET("/", mainController.Home)
	router.GET("/rule/*any", mainController.Error)
	router.POST("/validate", mainController.HandleValidationRequest)
	router.POST("/rule/:ruleName", mainController.HandleRuleRequest)

	// configure the app
	config.Configure()

	router.Run(":8080")
}

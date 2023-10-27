package main

import (
	"conectivity-checker-wizard/controllers"
	"conectivity-checker-wizard/rulemanager/builder"

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

	// build rules
	builder.BuildRules()

	// load all the templates
	router.LoadHTMLGlob("views/templates/*")

	// serve static files (e.g., css, js, images)
	router.Static("/views/static/", "./views/static/")

	// configure routes
	mainController := new(controllers.MainController)
	router.GET("/", mainController.Home)
	router.GET("/rule/*any", mainController.Error)
	router.GET("/cilium", mainController.CiliumPolicies)
	router.POST("/validate", mainController.HandleValidationRequest)
	router.POST("/rule/:ruleName", mainController.HandleOtherRequest)

	router.Run(":8080")
}

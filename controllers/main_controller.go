package controllers

import (
	"net/http"

	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/services"

	"github.com/gin-gonic/gin"
)

type MainController struct{}

func (mc *MainController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

// Binding request body (form data) with Object
func (mc *MainController) Execute(c *gin.Context) {
	var data models.FormData
	if err := c.ShouldBind(&data); err == nil {
		reponseTemplate := services.HandlRequest(c, data)
		c.HTML(http.StatusOK, reponseTemplate.Name, reponseTemplate)
	} else {
		// handle error page
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

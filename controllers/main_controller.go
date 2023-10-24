package controllers

import (
	"net/http"
	"strconv"

	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/services"
	"github.com/gin-gonic/gin"
)

type MainController struct{}

func (mc *MainController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

// func (mc *MainController) Questions(c *gin.Context) {
// 	questionID, _ := strconv.Atoi(c.Param("questionID"))
// 	services.HandlRequest(c, questionID)
// }

// Binding request body (form data) with Object
func (mc *MainController) Questions(c *gin.Context) {
	var data models.FormData
	questionID, _ := strconv.Atoi(c.Param("questionID"))
	if err := c.ShouldBind(&data); err == nil {
		services.HandlRequest(c, questionID, data)
	} else {
		// handle error page
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

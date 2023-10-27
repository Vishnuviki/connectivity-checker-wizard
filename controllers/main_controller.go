package controllers

import (
	"net/http"

	"conectivity-checker-wizard/cilium"
	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/services"

	"github.com/gin-gonic/gin"
)

type MainController struct{}

func (mc *MainController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func (mc *MainController) ValidateRule(c *gin.Context) {
	var data models.InputData
	if err := c.ShouldBind(&data); err == nil {
		reponseTemplate := services.HandleValidationRequest(c, data)
		c.HTML(http.StatusOK, reponseTemplate.Name, reponseTemplate)
	} else {
		// handle error page
		// TODO - Show Error or invalid page
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (mc *MainController) ExecuteRules(c *gin.Context) {
	if ruleName := c.Param("ruleName"); ruleName != "" {
		reponseTemplate := services.HandleRequest(c, ruleName)
		c.HTML(http.StatusOK, reponseTemplate.Name, reponseTemplate)
	} else {
		// handle error page
		// TODO - Show invalid page
	}
}

func (mc *MainController) CiliumPolicies(c *gin.Context) {
	policies, err := cilium.GetCiliumNetworkPolicies("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, policies)
}

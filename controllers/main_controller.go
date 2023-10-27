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

func (mc *MainController) ExecuteValidationRule(c *gin.Context) {
	var data models.InputData
	if err := c.ShouldBind(&data); err == nil {
		responseData := services.HandleValidationRequest(c, data)
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	} else {
		responseData := services.HandleInvalidRequest()
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	}
}

func (mc *MainController) ExecuteOtherRules(c *gin.Context) {
	if ruleName := c.Param("ruleName"); ruleName != "" {
		responseData := services.HandleRequest(c, ruleName)
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	} else {
		responseData := services.HandleInvalidRequest()
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	}
}

func (mc *MainController) Error(c *gin.Context) {
	responseData := services.HandleInvalidRequest()
	c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
}

func (mc *MainController) CiliumPolicies(c *gin.Context) {
	policies, err := cilium.GetCiliumNetworkPolicies("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, policies)
}

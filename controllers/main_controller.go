package controllers

import (
	"encoding/gob"
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/rulemanager/handler"
	"conectivity-checker-wizard/services/cilium"
	"conectivity-checker-wizard/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(models.InputData{})
}

type MainController struct{}

func (mc *MainController) Home(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	inputData := session.Get("inputData")
	session.Save()

	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"flashes":   flashes,
		"inputData": inputData,
	})
}

func (mc *MainController) HandleValidationRequest(c *gin.Context) {
	var data models.InputData
	session := sessions.Default(c)
	err := c.ShouldBind(&data)
	session.Set("inputData", data)
	session.Save()

	if err != nil {
		session.AddFlash("Validation failed: " + err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	} else if !data.IsDestinationAddressValid() {
		// TODO: what happen if user types `slack.com`??
		// We need to change the logic of checking FQDN
		session.AddFlash("Validation failed: Destination address not valid, must be IP or FQDN")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	} else {
		responseData := handler.HandleRules(c, constants.VALIDATION_RULE)
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	}
}

func (mc *MainController) HandleOtherRequest(c *gin.Context) {
	if ruleName := c.Param("ruleName"); ruleName != "" {
		responseData := handler.HandleRules(c, ruleName)
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	} else {
		responseData := handleInvalidRequest()
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	}
}

func (mc *MainController) Error(c *gin.Context) {
	responseData := handleInvalidRequest()
	c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
}

func (mc *MainController) CiliumPolicies(c *gin.Context) {
	policies, err := cilium.GetCiliumNetworkPolicies("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, policies)
}

func handleInvalidRequest() models.ResponseData {
	return utils.BuildResponseData(http.StatusNotFound, constants.PAGE_NOT_FOUND, "page-not-found.tmpl")
}

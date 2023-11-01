package controllers

import (
	"encoding/gob"
	"log"
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/rulemanager/executor"
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
	var inputData models.InputData
	session := sessions.Default(c)
	err := c.ShouldBind(&inputData)
	session.Set("inputData", inputData)
	session.Save()

	if err != nil {
		session.AddFlash("Validation failed: " + err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !inputData.IsDestinationAddressValid() {
		session.AddFlash("Validation failed: Destination address not valid, must be IP or FQDN")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !inputData.IsDestinationPortValid() {
		session.AddFlash("Validation failed: Invalid Destination Port")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	responseData := executor.ExecuteRules(c, constants.VALIDATION_RULE)
	c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
}

func (mc *MainController) HandleRuleRequest(c *gin.Context) {
	if ruleName := c.Param("ruleName"); ruleName != "" {
		responseData := executor.ExecuteRules(c, ruleName)
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	} else {
		responseData := utils.BuildInvalidResponseData()
		c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
	}
}

func (mc *MainController) Error(c *gin.Context) {
	log.Println("Invalid Request")
	responseData := utils.BuildInvalidResponseData()
	c.HTML(responseData.HTTPStatus, responseData.TemplateName, responseData)
}

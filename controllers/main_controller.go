package controllers

import (
	"net/http"

	"conectivity-checker-wizard/cilium"
	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/services"

	"encoding/gob"

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

func (mc *MainController) ValidateRule(c *gin.Context) {
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
	}

	if !data.IsDestinationAddressValid() {
		session.AddFlash("Validation failed: Destination address not valid, must be IP or FQDN")
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	reponseTemplate := services.HandleValidationRequest(c, data)
	c.HTML(http.StatusOK, reponseTemplate.Name, reponseTemplate)
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

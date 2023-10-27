package services

import (
	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func startSession(c *gin.Context, data models.InputData) {
	session := sessions.Default(c)
	session.Set("sourceNamespace", data.SourceNamespace)
	session.Set("destinationPort", data.DestinationPort)
	session.Set("destinationAddress", data.DestinationAddress)
	session.Save()
}

func HandleValidationRequest(c *gin.Context, formData models.InputData) models.ResponseData {
	session := sessions.Default(c)
	session.Clear()
	startSession(c, formData)

	if formData.IsDestinationAddressIP() {
		return HandleRules(c, "validationRule")
	} else {
		return HandleRules(c, "networkPolicyRule")
	}
}

func HandleRequest(c *gin.Context, ruleName string) models.ResponseData {
	return HandleRules(c, ruleName)
}

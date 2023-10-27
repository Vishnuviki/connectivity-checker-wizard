package services

import (
	"net/http"

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
	return HandleRules(c, "validationRule")
}

func HandleRequest(c *gin.Context, ruleName string) models.ResponseData {
	return HandleRules(c, ruleName)
}

func HandleInvalidRequest() models.ResponseData {
	return models.ResponseData{
		TemplateName: "page-not-found.tmpl",
		Content:      "Page Not Found.",
		HTTPStatus:   http.StatusNotFound,
	}
}

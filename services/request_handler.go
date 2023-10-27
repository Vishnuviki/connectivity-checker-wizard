package services

import (
	"net/http"

	"conectivity-checker-wizard/models"
	"github.com/gin-gonic/gin"
)

func HandleAllRules(c *gin.Context, ruleName string) models.ResponseData {
	return HandleRules(c, ruleName)
}

func HandleInvalidRequest() models.ResponseData {
	return models.ResponseData{
		TemplateName: "page-not-found.tmpl",
		Content:      "Page Not Found.",
		HTTPStatus:   http.StatusNotFound,
	}
}

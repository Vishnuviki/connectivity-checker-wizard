package services

import (
	"fmt"

	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func buildNoEgressPolicyResponse(c *gin.Context) models.TemplateData {
	templateData := models.TemplateData{}
	session := sessions.Default(c)
	destinationPort := session.Get("destinationPort")
	destinationAddress := session.Get("destinationAddress")
	templateData.Content = fmt.Sprintf("Oops, There is no network policy allowing this egress traffic"+
		"link-to-docs-about-egress-policy - "+
		"destinationPort: %s, destinationAddress: %s", destinationPort, destinationAddress)
	session.Clear()
	templateData.Name = "final-response.tmpl"
	return templateData
}

func buildDefaultResponse() models.TemplateData {
	templateData := models.TemplateData{}
	templateData.Name = "thankyou.html"
	return templateData
}

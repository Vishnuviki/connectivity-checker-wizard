package services

import (
	"fmt"

	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func buildQuestionOne(c *gin.Context) models.TemplateData {
	templateData := models.TemplateData{}
	session := sessions.Default(c)
	destinationAddress := session.Get("destinationAddress")
	templateData.Content = fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", destinationAddress)
	templateData.Name = "final-question.tmpl"
	templateData.HTTPMethod = "post"
	templateData.Endpoint = "/api/question/2"
	return templateData
}

func buildQuestionTwo(c *gin.Context) models.TemplateData {
	templateData := models.TemplateData{}
	session := sessions.Default(c)
	sourceNamespace := session.Get("sourceNamespace")
	templateData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Now we will test the DNS lookup", sourceNamespace)
	templateData.Name = "final-question.tmpl"
	templateData.HTTPMethod = "get"
	templateData.Endpoint = "/api/condition/2"
	return templateData
}

func buildQuestionThree(c *gin.Context) models.TemplateData {
	templateData := models.TemplateData{}
	session := sessions.Default(c)
	sourceNamespace := session.Get("sourceNamespace")
	templateData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", sourceNamespace)
	templateData.Name = "final-question.tmpl"
	templateData.HTTPMethod = "get"
	templateData.Endpoint = "/api/condition/3"
	return templateData
}

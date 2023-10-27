package rules

import (
	"conectivity-checker-wizard/models"

	// "github.com/gin-contrib/sessions"
)

func BuildResponseData(httpStatus int, content, templateName string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = content
	responseData.TemplateName = templateName
	responseData.HTTPStatus = httpStatus
	return *responseData
}

// func buildInputData(session sessions.Session) models.InputData {
// 	sourceNamespace := session.Get("sourceNamespace").(string)
// 	destinationPort := session.Get("destinationPort").(string)
// 	destinationAddress := session.Get("destinationAddress").(string)
// 	return *models.NewInputData(sourceNamespace, destinationPort, destinationAddress)
// }

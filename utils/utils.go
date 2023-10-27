package utils

import (
	"conectivity-checker-wizard/models"
)

func BuildResponseData(httpStatus int, content, templateName string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = content
	responseData.TemplateName = templateName
	responseData.HTTPStatus = httpStatus
	return *responseData
}

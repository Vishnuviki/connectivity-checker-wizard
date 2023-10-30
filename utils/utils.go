package utils

import (
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
)

/*
	Add reuseable code here
*/

func BuildInvalidResponseData() models.ResponseData {
	return models.ResponseData{
		HTTPStatus:      http.StatusNotFound,
		TemplateName:    "page-not-found.tmpl",
		TemplateContent: constants.PAGE_NOT_FOUND_MESSAGE,
	}
}

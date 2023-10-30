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
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusNotFound).
		WithTemplateName("page-not-found.tmpl").
		WithTemplateContent(constants.PAGE_NOT_FOUND_MESSAGE).
		Build()
}

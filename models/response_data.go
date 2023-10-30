package models

type ResponseData struct {
	HTTPStatus         int
	TemplateName       string
	TemplateContent    string
	TemplateFormAction string
	TemplateFormMethod string
}

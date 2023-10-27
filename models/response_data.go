package models

type ResponseData struct {
	HTTPStatus   int
	HTTPMethod   string
	TemplateName string
	Content      string
	Endpoint     string
}

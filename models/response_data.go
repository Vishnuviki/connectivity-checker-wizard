package models

type ResponseData struct {
	HTTPStatus         int
	TemplateName       string
	TemplateContent    string
	TemplateFormAction string
	TemplateFormMethod string
}

type ResponseDataBuilder struct {
	responseData ResponseData
}

func NewResponseDataBuilder() *ResponseDataBuilder {
	return new(ResponseDataBuilder)
}

func (rb *ResponseDataBuilder) WithHTTPStatus(httStatus int) *ResponseDataBuilder {
	rb.responseData.HTTPStatus = httStatus
	return rb
}

func (rb *ResponseDataBuilder) WithTemplateName(templateName string) *ResponseDataBuilder {
	rb.responseData.TemplateName = templateName
	return rb
}

func (rb *ResponseDataBuilder) WithTemplateContent(templateContent string) *ResponseDataBuilder {
	rb.responseData.TemplateContent = templateContent
	return rb
}

func (rb *ResponseDataBuilder) WithTemplateFormAction(templateFormAction string) *ResponseDataBuilder {
	rb.responseData.TemplateFormAction = templateFormAction
	return rb
}

func (rb *ResponseDataBuilder) WithTemplateFormMethod(templateFormMethod string) *ResponseDataBuilder {
	rb.responseData.TemplateFormMethod = templateFormMethod
	return rb
}

func (rb *ResponseDataBuilder) Build() ResponseData {
	return rb.responseData
}

package models

type ResponseData struct {
	HTTPStatus   int
	HTTPMethod   string
	TemplateName string
	Content      string
	Endpoint     string
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

func (rb *ResponseDataBuilder) WithHTTPMethod(httMethod string) *ResponseDataBuilder {
	rb.responseData.HTTPMethod = httMethod
	return rb
}

func (rb *ResponseDataBuilder) WithTemplateName(templateName string) *ResponseDataBuilder {
	rb.responseData.TemplateName = templateName
	return rb
}

func (rb *ResponseDataBuilder) WithContent(content string) *ResponseDataBuilder {
	rb.responseData.Content = content
	return rb
}

func (rb *ResponseDataBuilder) WithEndpoint(endpoint string) *ResponseDataBuilder {
	rb.responseData.Endpoint = endpoint
	return rb
}

func (rb *ResponseDataBuilder) Build() ResponseData {
	return rb.responseData
}

package models

import "testing"

func TestResponseDataBuilder(t *testing.T) {
	responseData := NewResponseDataBuilder().
		WithHTTPStatus(200).
		WithTemplateName("home.tmpl").
		WithTemplateContent("Page not found").
		WithTemplateFormAction("/rule/validationRule").
		WithTemplateFormMethod("POST").
		Build()

	if responseData.HTTPStatus != 200 {
		t.Errorf("Expected HTTPStatus to be 200, got %d", responseData.HTTPStatus)
	}

	if responseData.TemplateFormMethod != "POST" {
		t.Errorf("Expected HTTPMethod to be 'POST', got %s", responseData.TemplateFormMethod)
	}

	if responseData.TemplateName != "home.tmpl" {
		t.Errorf("Expected TemplateName to be 'home.tmpl', got %s", responseData.TemplateName)
	}

	if responseData.TemplateContent != "Page not found" {
		t.Errorf("Expected Content to be 'Some content', got %s", responseData.TemplateContent)
	}

	if responseData.TemplateFormAction != "/rule/validationRule" {
		t.Errorf("Expected Endpoint to be '/rule/validationRule', got %s", responseData.TemplateFormAction)
	}
}

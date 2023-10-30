package rules

import (
	"fmt"
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
)

type DNSLookUPRule struct {
	name     string
	nextRule i.Rule
}

func (r *DNSLookUPRule) SetNextRule(nextRule i.Rule) {
	r.nextRule = nextRule
}

func (r *DNSLookUPRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *DNSLookUPRule) Execute(inputData models.InputData) models.ResponseData {
	content := fmt.Sprintf("This is a %s Page", constants.DNS_LOOK_UP_RULE)
	return models.ResponseData{
		HTTPStatus:      http.StatusOK,
		TemplateName:    "response.tmpl",
		TemplateContent: content,
	}
}

package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
)

type DispatchIPRule struct {
	name     string
	nextRule i.Rule
}

func (r *DispatchIPRule) SetNextRule(nextRule i.Rule) {
	r.nextRule = nextRule
}

func (r *DispatchIPRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *DispatchIPRule) Execute(inputData models.InputData) models.ResponseData {
	log.Printf("Executing Rule: %s", r.name)
	content := fmt.Sprintf("This is a %s Page", constants.DISPATCH_IP_RULE)
	return models.ResponseData{
		HTTPStatus:      http.StatusOK,
		TemplateName:    "response.tmpl",
		TemplateContent: content,
	}
}

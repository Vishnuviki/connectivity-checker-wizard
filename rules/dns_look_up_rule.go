package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/models"
)

type DNSLookUPRule struct {
	name     string
	nextRule Rule
}

func (r *DNSLookUPRule) SetNextRule(nextRule Rule) {
	r.nextRule = nextRule
}

func (r *DNSLookUPRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *DNSLookUPRule) Execute(inputData models.InputData) models.ResponseData {
	log.Printf("Executing Rule: %s", DNS_LOOK_UP_RULE)
	if r.nextRule != nil {
		r.nextRule.Execute(inputData)
	}
	content := fmt.Sprintf("This is a %s Page", DNS_LOOK_UP_RULE)
	return BuildResponseData(http.StatusOK, content, "response.tmpl")
}

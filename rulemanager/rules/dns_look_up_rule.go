package rules

import (
	"fmt"
	"log"
	"net/http"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/utils"
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
	log.Printf("Executing Rule: %s", c.DNS_LOOK_UP_RULE)
	if r.nextRule != nil {
		r.nextRule.Execute(inputData)
	}
	content := fmt.Sprintf("This is a %s Page", c.DNS_LOOK_UP_RULE)
	return utils.BuildResponseData(http.StatusOK, content, "response.tmpl")
}

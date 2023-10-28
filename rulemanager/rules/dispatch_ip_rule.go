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
	log.Printf("Executing Rule: %s", c.DISPATCH_IP_RULE)
	content := fmt.Sprintf("This is a %s Page", c.DISPATCH_IP_RULE)
	return utils.BuildResponseData(http.StatusOK, content, "response.tmpl")
}

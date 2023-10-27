package rules

import (
	"fmt"
	"log"
	"net/http"

	i "conectivity-checker-wizard/interfaces"
	"conectivity-checker-wizard/models"
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
	log.Printf("Executing Rule: %s", DISPATCH_IP_RULE)
	content := fmt.Sprintf("This is a %s Page", DISPATCH_IP_RULE)
	return utils.BuildResponseData(http.StatusOK, content, "response.tmpl")
}

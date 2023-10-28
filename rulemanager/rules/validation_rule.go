package rules

import (
	"fmt"
	"log"
	"net/http"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
)

type ValidationRule struct {
	name     string
	nextRule i.Rule
}

func (r *ValidationRule) SetNextRule(nextRule i.Rule) {
	r.nextRule = nextRule
}

func (r *ValidationRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *ValidationRule) Execute(inputData models.InputData) models.ResponseData {
	log.Printf("Executing Rule: %s", c.VALIDATION_RULE)
	if inputData.IsDestinationAddressIP() {
		return buildResponse(inputData.DestinationAddress)
	} else {
		// execute networkPolicyRule
		return r.nextRule.Execute(inputData)
	}
}

func buildResponse(destinationAddress string) models.ResponseData {
	content := fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", destinationAddress)
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithHTTPMethod("post").
		WithTemplateName("question.tmpl").
		WithEndpoint("/rule/networkPolicyRule").
		WithContent(content).
		Build()
}

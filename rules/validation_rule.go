package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ValidationRule struct {
	name     string
	nextRule Rule // may be we can make it a List
}

func (r *ValidationRule) SetNextRule(nextRule Rule) {
	r.nextRule = nextRule
}

func (r *ValidationRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *ValidationRule) Execute(c *gin.Context) models.ResponseData {
	log.Printf("Executing Rule: %s", VALIDATION_RULE)
	session := sessions.Default(c)
	inputData := buildInputData(session)
	if inputData.IsDestinationAddressIP() {
		return buildResponse(inputData.DestinationAddress)
	} else {
		// execute networkPolicyRule
		return r.nextRule.Execute(c)
	}
}

func buildResponse(address string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", address)
	responseData.TemplateName = "question.tmpl"
	responseData.HTTPMethod = "post"
	responseData.HTTPStatus = http.StatusOK
	responseData.Endpoint = "/rule/networkPolicyRule"
	return *responseData
}

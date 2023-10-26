package rules

import (
	"log"

	"conectivity-checker-wizard/models"

	"github.com/gin-gonic/gin"
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

func (r *DNSLookUPRule) Execute(c *gin.Context) models.ResponseData {
	log.Printf("Executing Rule: %s", DNS_LOOK_UP_RULE)
	if r.nextRule != nil {
		r.nextRule.Execute(c)
	}
	return buildDefaultResponse(DNS_LOOK_UP_RULE)
}

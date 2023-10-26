package rules

import (
	"log"

	"conectivity-checker-wizard/models"

	"github.com/gin-gonic/gin"
)

type DispatchIPRule struct {
	name     string
	nextRule Rule
}

func (r *DispatchIPRule) SetNextRule(nextRule Rule) {
	r.nextRule = nextRule
}

func (r *DispatchIPRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *DispatchIPRule) Execute(c *gin.Context) models.ResponseData {
	log.Printf("Executing Rule: %s", DISPATCH_IP_RULE)
	return buildDefaultResponse(DISPATCH_IP_RULE)
}

package rules

import (
	"fmt"
	"log"
	"net/http"

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
	content := fmt.Sprintf("This is a %s Page", DISPATCH_IP_RULE)
	return BuildResponseData(http.StatusOK, content, "response.tmpl")
}

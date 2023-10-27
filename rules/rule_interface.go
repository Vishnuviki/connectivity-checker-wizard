package rules

import (
	"conectivity-checker-wizard/models"
	"github.com/gin-gonic/gin"
)

type Rule interface {
	SetNextRule(nextRule Rule)
	SetName(ruleName string)
	Execute(c *gin.Context) models.ResponseData
}

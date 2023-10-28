package handler

import (
	"net/http"

	"conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	"conectivity-checker-wizard/rulemanager/builder"
	"conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var ruleMap = rulemap.GetInstance()

func BuildRuleMap() {
	builder.BuildRules(ruleMap)
}

func HandleRules(c *gin.Context, ruleName string) models.ResponseData {
	if rule, ok := ruleMap.GetRuleByName(ruleName); ok {
		session := sessions.Default(c)
		inputData := session.Get("inputData").(models.InputData)
		// execute rule
		return rule.Execute(inputData)
	} else {
		return utils.BuildResponseData(http.StatusNotFound, constants.PAGE_NOT_FOUND, "page-not-found.tmpl")
	}
}

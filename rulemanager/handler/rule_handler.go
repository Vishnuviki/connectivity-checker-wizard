package handler

import (
	"net/http"

	"conectivity-checker-wizard/models"
	rm "conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandleRules(c *gin.Context, ruleName string) models.ResponseData {
	ruleMap := rm.GetInstance()
	if rule, ok := ruleMap.GetRuleByName(ruleName); ok {
		session := sessions.Default(c)
		inputData := session.Get("inputData").(models.InputData)
		// execute rule
		return rule.Execute(inputData)
	} else {
		return utils.BuildResponseData(http.StatusNotFound, "Page Not Found.", "page-not-found.tmpl")
	}
}

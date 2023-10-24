package services

import (
	"strconv"

	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandlRequest(c *gin.Context, data models.FormData) models.TemplateData {
	var responseData models.TemplateData
	questionID, _ := strconv.Atoi(c.Param("questionID"))

	switch session := sessions.Default(c); questionID {
	case 1:
		session.Clear()
		startSession(c, data)
		responseData = buildQuestionOne(c)
	case 2:
		// API Call to API-SERVER - NW policy allows egress and destination is a hostname
		sourceNamespace := session.Get("sourceNamespace").(string)
		res := getCiliumNetworkPolicy(sourceNamespace)
		if res.IsHostname {
			responseData = buildQuestionTwo(c)
		} else if res.IsIPAddress {
			responseData = buildQuestionThree(c)
		} else {
			return buildNoEgressPolicyResponse(c)
		}
	default:
		session.Clear()
		responseData = buildDefaultResponse()
	}
	return responseData
}

func startSession(c *gin.Context, data models.FormData) {
	session := sessions.Default(c)
	session.Set("sourceNamespace", data.SourceNamespace)
	session.Set("destinationPort", data.DestinationPort)
	session.Set("destinationAddress", data.DestinationAddress)
	session.Save()
}

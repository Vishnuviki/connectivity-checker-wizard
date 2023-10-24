package services

import (
	"conectivity-checker-wizard/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlRequest(c *gin.Context, questionID int, data models.FormData) {
	switch session := sessions.Default(c); questionID {
	case 1:
		session.Clear()
		startSession(c, data)
		buildQuestionOne(c)
	case 2:
		// API Call to API-SERVER - NW policy allows egress and destination is a hostname
		sourceNamespace := session.Get("sourceNamespace").(string)
		res := getCiliumNetworkPolicy(sourceNamespace)
		if res.IsHostname {
			buildQuestionTwo(c)
		} else if res.IsIPAddress {
			buildQuestionThree(c)
		} else {
			buildNoEgressPolicyResponse(c)
		}
	default:
		session.Clear()
		c.HTML(http.StatusOK, "thankyou.html", nil)
	}
}

func startSession(c *gin.Context, data models.FormData) {
	session := sessions.Default(c)
	// Parsing values from the Form
	// sourceCluster := c.PostForm("sourceCluster")
	// sourceNamespace := c.PostForm("sourceNamespace")
	// destinationPort := c.PostForm("destinationPort")
	// destinationAddress := c.PostForm("destinationAddress")

	// setting values to Session
	// session.Set("sourceCluster", sourceCluster)
	// session.Set("sourceNamespace", sourceNamespace)
	// session.Set("destinationPort", destinationPort)
	// session.Set("destinationAddress", destinationAddress)

	session.Set("sourceNamespace", data.SourceNamespace)
	session.Set("destinationPort", data.DestinationPort)
	session.Set("destinationAddress", data.DestinationAddress)
	session.Save()
}

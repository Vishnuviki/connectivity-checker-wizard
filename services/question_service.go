package services

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func buildQuestionOne(c *gin.Context) {
	session := sessions.Default(c)
	destinationAddress := session.Get("destinationAddress")
	fmt.Println("DESTINATION-ADDRESS: ", destinationAddress)
	question := fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", destinationAddress)
	c.HTML(http.StatusOK, "final-question.tmpl", gin.H{
		"question":   question,
		"httpMethod": "post",
		"endpoint":   "/api/question/2",
	})
}

func buildQuestionTwo(c *gin.Context) {
	session := sessions.Default(c)
	sourceNamespace := session.Get("sourceNamespace")
	question := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Now we will test the DNS lookup", sourceNamespace)
	c.HTML(http.StatusOK, "final-question.tmpl", gin.H{
		"question":   question,
		"httpMethod": "get",
		"endpoint":   "/api/look-up-target",
	})
}

func buildQuestionThree(c *gin.Context) {
	session := sessions.Default(c)
	sourceNamespace := session.Get("sourceNamespace")
	question := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", sourceNamespace)
	c.HTML(http.StatusOK, "final-question.tmpl", gin.H{
		"question":   question,
		"httpMethod": "get",
		"endpoint":   "/api/dns-lookup",
	})
}

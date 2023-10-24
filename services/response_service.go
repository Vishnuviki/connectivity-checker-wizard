package services

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func buildNoEgressPolicyResponse(c *gin.Context) {
	session := sessions.Default(c)
	destinationPort := session.Get("destinationPort")
	destinationAddress := session.Get("destinationAddress")
	response := fmt.Sprintf("Oops, There is no network policy allowing this egress traffic"+
		"link-to-docs-about-egress-policy - "+
		"destinationPort: %s, destinationAddress: %s", destinationPort, destinationAddress)
	session.Clear()
	c.HTML(http.StatusOK, "final-response.tmpl", gin.H{
		"response": response,
	})
}

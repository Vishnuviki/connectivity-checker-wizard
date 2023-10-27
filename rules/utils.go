package rules

import (
	"fmt"

	"conectivity-checker-wizard/fsm"
	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
)

func buildDefaultResponse(ruleName string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("This is a %s page", ruleName)
	responseData.Name = "response.tmpl"
	return *responseData
}

func buildInputData(session sessions.Session) fsm.InputData {
	sourceNamespace := session.Get("sourceNamespace").(string)
	destinationPort := session.Get("destinationPort").(string)
	destinationAddress := session.Get("destinationAddress").(string)
	return *fsm.NewInputData(sourceNamespace, destinationPort, destinationAddress)
}

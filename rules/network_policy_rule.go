package rules

import (
	"fmt"
	"log"

	"conectivity-checker-wizard/fsm"
	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NetworkPolicyRule struct {
	name     string
	nextRule Rule
}

func (r *NetworkPolicyRule) SetNextRule(nextRule Rule) {
	r.nextRule = nextRule
}

func (r *NetworkPolicyRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *NetworkPolicyRule) Execute(c *gin.Context) models.ResponseData {
	log.Printf("Executing Rule: %s", NETWORK_POLICY_RULE)
	session := sessions.Default(c)
	inputData := buildInputData(session)

	if inputData.IsDestinationAddressIP() {
		fmt.Println("IsDestinationAddressIP")
		return processIPAddressRequest(inputData)
	} else {
		fmt.Println("IsDestinationAddressFQDN")
		return processFQDNRequest(inputData)
	}
}

func processFQDNRequest(input fsm.InputData) models.ResponseData {
	isAvailable := checkFQDNEgreesRules(input.SourceNamespace, input.DestinationAddress)
	if isAvailable {
		return buildFQDNResponse(input.SourceNamespace)
	} else {
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func processIPAddressRequest(input fsm.InputData) models.ResponseData {
	isAvailable := checkIPAddressEgressRules(input.SourceNamespace, input.DestinationAddress)
	if isAvailable {
		return buildIPAddressResponse(input.SourceNamespace)
	} else {
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func buildFQDNResponse(namespace string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Now we will test the DNS lookup", namespace)
	responseData.Name = "question.tmpl"
	responseData.HTTPMethod = "post"
	responseData.Endpoint = "/rule/dnsLookUPRule"
	return *responseData
}

func buildIPAddressResponse(namespace string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", namespace)
	responseData.Name = "question.tmpl"
	responseData.HTTPMethod = "post"
	responseData.Endpoint = "/rule/dispatchIPRule"
	return *responseData
}

func buildNoEgressPolicyResponse(port, address string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Oops, There is no network policy allowing this egress traffic - "+
		"link-to-docs-about-egress-policy - "+
		"destinationPort: %s, destinationAddress: %s", port, address)
	responseData.Name = "response.tmpl"
	return *responseData
}

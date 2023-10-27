package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/cilium"
	"conectivity-checker-wizard/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NetworkPolicyRule struct {
	name                string
	nextRule            Rule
	ciliumPolicyChecker cilium.CiliumPolicyChecker
}

func NewNetworkPolicyRule(name string, nextRule Rule, ciliumPolicyChecker cilium.CiliumPolicyChecker) *NetworkPolicyRule {
	return &NetworkPolicyRule{
		name:                name,
		nextRule:            nextRule,
		ciliumPolicyChecker: ciliumPolicyChecker,
	}
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
		return r.processIPAddressRequest(inputData)
	} else {
		fmt.Println("IsDestinationAddressFQDN")
		return r.processFQDNRequest(inputData)
	}
}

func (r *NetworkPolicyRule) processFQDNRequest(input models.InputData) models.ResponseData {
	isAvailable, err := r.ciliumPolicyChecker.CheckFQDNAllowedByPolicyInNamespace(input.DestinationAddress, input.SourceNamespace)
	if err != nil {
		// TODO: error handling ??? we cannot just ignore it or say it is not allowed.
		// true/false from CheckFQDNAllowedByPolicyInNamespace is not enough
		// if there is an error we probably need to surface it to the user?
		log.Println(err)
	}
	if isAvailable {
		return buildFQDNResponse(input.SourceNamespace)
	} else {
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func (r *NetworkPolicyRule) processIPAddressRequest(input models.InputData) models.ResponseData {
	isAvailable, err := r.ciliumPolicyChecker.CheckIPAllowedByPolicyInNamespace(input.DestinationAddress, input.SourceNamespace)
	if err != nil {
		// TODO: error handling ??? we cannot just ignore it or say it is not allowed.
		// true/false from CheckIPAllowedByPolicyInNamespace is not enough
		// if there is an error we probably need to surface it to the user?
		log.Println(err)
	}
	if isAvailable {
		return buildIPAddressResponse(input.SourceNamespace)
	} else {
		fmt.Println("ELSE")
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func buildFQDNResponse(namespace string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Now we will test the DNS lookup", namespace)
	responseData.TemplateName = "question.tmpl"
	responseData.HTTPMethod = "post"
	responseData.HTTPStatus = http.StatusOK
	responseData.Endpoint = "/rule/dnsLookUPRule"
	return *responseData
}

func buildIPAddressResponse(namespace string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", namespace)
	responseData.TemplateName = "question.tmpl"
	responseData.HTTPMethod = "post"
	responseData.HTTPStatus = http.StatusOK
	responseData.Endpoint = "/rule/dispatchIPRule"
	return *responseData
}

func buildNoEgressPolicyResponse(port, address string) models.ResponseData {
	responseData := new(models.ResponseData)
	responseData.Content = fmt.Sprintf("Oops, There is no network policy allowing this egress traffic - "+
		"link-to-docs-about-egress-policy - "+
		"destinationPort: %s, destinationAddress: %s", port, address)
	responseData.TemplateName = "response.tmpl"
	responseData.HTTPStatus = http.StatusOK
	return *responseData
}

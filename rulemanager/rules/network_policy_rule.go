package rules

import (
	"fmt"
	"log"
	"net/http"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/services/cilium"
)

type NetworkPolicyRule struct {
	name                string
	nextRule            i.Rule
	ciliumPolicyChecker cilium.CiliumPolicyChecker
}

func NewNetworkPolicyRule(name string, nextRule i.Rule, ciliumPolicyChecker cilium.CiliumPolicyChecker) *NetworkPolicyRule {
	return &NetworkPolicyRule{
		name:                name,
		nextRule:            nextRule,
		ciliumPolicyChecker: ciliumPolicyChecker,
	}
}

func (r *NetworkPolicyRule) SetNextRule(nextRule i.Rule) {
	r.nextRule = nextRule
}

func (r *NetworkPolicyRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *NetworkPolicyRule) Execute(inputData models.InputData) models.ResponseData {
	log.Printf("Executing Rule: %s", c.NETWORK_POLICY_RULE)
	if inputData.IsDestinationAddressIP() {
		return r.processIPAddressRequest(inputData)
	} else {
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
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func buildFQDNResponse(sourceNamespace string) models.ResponseData {
	content := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Now we will test the DNS lookup", sourceNamespace)
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithHTTPMethod("post").
		WithTemplateName("question.tmpl").
		WithEndpoint("/rule/dnsLookUPRule").
		WithContent(content).
		Build()
}

func buildIPAddressResponse(sourceNamespace string) models.ResponseData {
	content := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", sourceNamespace)
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithHTTPMethod("post").
		WithTemplateName("question.tmpl").
		WithEndpoint("/rule/dispatchIPRule").
		WithContent(content).
		Build()
}

func buildNoEgressPolicyResponse(port, address string) models.ResponseData {
	// TODO - Append with relevant core-docs link
	content := fmt.Sprintf("Oops, There is no network policy allowing this egress traffic - " +
		"link-to-docs-about-egress-policy")
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithTemplateName("response.tmpl").
		WithContent(content).
		Build()
}

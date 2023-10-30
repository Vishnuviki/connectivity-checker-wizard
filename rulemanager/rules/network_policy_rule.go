package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/services/cilium"
)

type NetworkPolicyRule struct {
	name     string
	nextRule i.Rule
}

func (r *NetworkPolicyRule) SetName(ruleName string) {
	r.name = ruleName
}

func (r *NetworkPolicyRule) SetNextRule(nextRule i.Rule) {
	r.nextRule = nextRule
}

func (r *NetworkPolicyRule) Execute(inputData models.InputData) models.ResponseData {
	log.Printf("Executing Rule: %s", r.name)
	if inputData.IsDestinationAddressIP() {
		return processIPAddressRequest(inputData)
	} else {
		return processFQDNRequest(inputData)
	}
}

func processFQDNRequest(input models.InputData) models.ResponseData {
	ciliumClient := cilium.NewCiliumClient()
	policyList, err := ciliumClient.GetCiliumNetworkPolicies(input.SourceNamespace)
	if err != nil {
		return buildErrorResponse()
	}
	isAvailable := ciliumClient.CheckFQDNAllowedByPolicyInNamespace(input.DestinationAddress, policyList)
	if isAvailable {
		return buildFQDNResponse(input.SourceNamespace)
	} else {
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func processIPAddressRequest(input models.InputData) models.ResponseData {
	ciliumClient := cilium.NewCiliumClient()
	policyList, err := ciliumClient.GetCiliumNetworkPolicies(input.SourceNamespace)
	if err != nil {
		return buildErrorResponse()
	}
	isAvailable := ciliumClient.CheckIPAllowedByPolicyInNamespace(input.DestinationAddress, policyList)
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
		WithTemplateName("question.tmpl").
		WithTemplateContent(content).
		WithTemplateFormMethod(http.MethodPost).
		WithTemplateFormAction("/rule/dnsLookUPRule").
		Build()
}

func buildIPAddressResponse(sourceNamespace string) models.ResponseData {
	content := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", sourceNamespace)
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithTemplateName("question.tmpl").
		WithTemplateContent(content).
		WithTemplateFormMethod(http.MethodPost).
		WithTemplateFormAction("/rule/dispatchIPRule").
		Build()
}

func buildNoEgressPolicyResponse(port, address string) models.ResponseData {
	// TODO - Append with relevant core-docs link
	content := fmt.Sprintf("Oops, There is no network policy allowing this egress traffic - " +
		"link-to-docs-about-egress-policy")
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithTemplateName("response.tmpl").
		WithTemplateContent(content).
		Build()
}

func buildErrorResponse() models.ResponseData {
	content := "We apologize for the inconvenience, as we're currently encountering some technical issues. " +
		"Please get in touch with #core-support channel for further assistance."
	return models.NewResponseDataBuilder().
		WithHTTPStatus(http.StatusOK).
		WithTemplateName("response.tmpl").
		WithTemplateContent(content).
		Build()
}

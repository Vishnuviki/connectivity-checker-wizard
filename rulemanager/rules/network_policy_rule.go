package rules

import (
	"fmt"
	"log"
	"net/http"

	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/cilium"
)

type NetworkPolicyRule struct {
	name     string
	nextRule i.Rule
	policyChecker cilium.PolicyChecker
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
		return r.processIPAddressRequest(inputData)
	} else {
		return r.processFQDNRequest(inputData)
	}
}

func (r *NetworkPolicyRule) processFQDNRequest(input models.InputData) models.ResponseData {
	isAvailable, err := r.policyChecker.CheckFQDNAllowedByPolicyInNamespace(input.DestinationAddress, input.SourceNamespace)
	if err != nil {
		// TODO: error handling ??? we cannot just ignore it or say it is not allowed.
		log.Println(err)
	}
	if isAvailable {
		return buildFQDNResponse(input.SourceNamespace)
	} else {
		return buildNoEgressPolicyResponse(input.DestinationPort, input.DestinationAddress)
	}
}

func (r *NetworkPolicyRule) processIPAddressRequest(input models.InputData) models.ResponseData {
	isAvailable, err := r.policyChecker.CheckIPAllowedByPolicyInNamespace(input.DestinationAddress, input.SourceNamespace)
	if err != nil {
		// TODO: error handling ??? we cannot just ignore it or say it is not allowed.
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
	return models.ResponseData{
		HTTPStatus:         http.StatusOK,
		TemplateName:       "question.tmpl",
		TemplateContent:    content,
		TemplateFormMethod: http.MethodPost,
		TemplateFormAction: "/rule/dnsLookUPRule",
	}
}

func buildIPAddressResponse(sourceNamespace string) models.ResponseData {
	content := fmt.Sprintf("Good news, the source Namespace (%v) has a network policy allowing this traffic out. "+
		"Because the destination is an IP address, we don't need to examine DNS", sourceNamespace)
	return models.ResponseData{
		HTTPStatus:         http.StatusOK,
		TemplateName:       "question.tmpl",
		TemplateContent:    content,
		TemplateFormMethod: http.MethodPost,
		TemplateFormAction: "/rule/dispatchIPRule",
	}
}

func buildNoEgressPolicyResponse(port, address string) models.ResponseData {
	// TODO - Append with relevant core-docs link
	content := fmt.Sprintf("Oops, There is no network policy allowing this egress traffic - " +
		"link-to-docs-about-egress-policy")
	return models.ResponseData{
		HTTPStatus:      http.StatusOK,
		TemplateName:    "response.tmpl",
		TemplateContent: content,
	}
}

func buildErrorResponse() models.ResponseData {
	content := "We apologize for the inconvenience, as we're currently encountering some technical issues. " +
		"Please get in touch with #core-support channel for further assistance."
	return models.ResponseData{
		HTTPStatus:      http.StatusOK,
		TemplateName:    "response.tmpl",
		TemplateContent: content,
	}
}

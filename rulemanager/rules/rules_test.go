package rules

import (
	"fmt"
	"testing"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/models"
	i "conectivity-checker-wizard/rulemanager/interfaces"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAllRules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Rules Suite")
}

var _ = Describe("Testing Rules", func() {

	Context("Test Validation Rule", func() {
		var (
			validationRule i.Rule
			inputData      *models.InputData
		)

		BeforeEach(func() {
			validationRule = buildValidationRule()
			inputData = models.NewInputData("default", "8080", "192.168.1.0")
		})

		It("Should call buildResponse() function and returns responseData", func() {
			responseData := validationRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("question.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildValidationRuleResponseContent()))
			Expect(responseData.TemplateFormMethod).To(Equal("POST"))
			Expect(responseData.TemplateFormAction).To(Equal("/rule/networkPolicyRule"))
		})

		// TODO: Improve the test by using Mock Cilium client and return success response.
		It("Should call nextRule Execute() method if the destinationAddress is a FQDN", func() {
			inputData.DestinationAddress = "sky.slack.com"
			responseData := validationRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("response.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildErrorResponseContent()))
		})

		It("Should return a ResponseData with a message related to destinationAddress", func() {
			responseData := buildResponse(inputData.DestinationAddress)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("question.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildValidationRuleResponseContent()))
			Expect(responseData.TemplateFormMethod).To(Equal("POST"))
			Expect(responseData.TemplateFormAction).To(Equal("/rule/networkPolicyRule"))
		})
	})

	Context("Test DNSLookUPRule Rule", func() {
		var (
			dnsLookUPRule i.Rule
			inputData     *models.InputData
		)

		BeforeEach(func() {
			dnsLookUPRule = buildDNSLookUPRule()
			inputData = models.NewInputData("default", "8080", "192.168.1.0")
		})

		It("Should execute and return a ResponseData with the expected content", func() {
			responseData := dnsLookUPRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("response.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildDNSLookUPRuleResponseContent()))
		})
	})

	Context("Test DispatchIPRule Rule", func() {
		var (
			dispatchIPRule i.Rule
			inputData      *models.InputData
		)

		BeforeEach(func() {
			dispatchIPRule = buildDispatchIPRule()
			inputData = models.NewInputData("default", "8080", "192.168.1.0")
		})

		It("Should execute and return a ResponseData with the expected content", func() {
			responseData := dispatchIPRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("response.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildDispatchIPRuleResponseContent()))
		})
	})

})

func buildValidationRule() i.Rule {
	validationRule := new(ValidationRule)
	validationRule.SetName(c.VALIDATION_RULE)
	validationRule.SetNextRule(new(NetworkPolicyRule))
	return validationRule
}

func buildValidationRuleResponseContent() string {
	return fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", "192.168.1.0")
}

func buildErrorResponseContent() string {
	return "We apologize for the inconvenience, as we're currently encountering some technical issues. " +
		"Please get in touch with #core-support channel for further assistance."
}

func buildDNSLookUPRule() i.Rule {
	dnsLookUPRule := new(DNSLookUPRule)
	dnsLookUPRule.SetName(c.DNS_LOOK_UP_RULE)
	dnsLookUPRule.SetNextRule(nil)
	return dnsLookUPRule
}

func buildDNSLookUPRuleResponseContent() string {
	return "This is a dnsLookUPRule Page"
}

func buildDispatchIPRule() i.Rule {
	dispatchIPRule := new(DispatchIPRule)
	dispatchIPRule.SetName(c.DISPATCH_IP_RULE)
	dispatchIPRule.SetNextRule(nil)
	return dispatchIPRule
}

func buildDispatchIPRuleResponseContent() string {
	return "This is a dispatchIPRule Page"
}

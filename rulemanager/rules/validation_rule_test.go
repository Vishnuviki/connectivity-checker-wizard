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

func TestValidationRule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test ValidationRule Suite")
}

var _ = Describe("Testing ValidationRule functions", func() {
	var (
		validationRule i.Rule
		inputData      *models.InputData
	)

	BeforeEach(func() {
		validationRule = buildValidationRule()
		inputData = models.NewInputData("default", "8080", "192.168.1.0")
	})

	Context("Test Execute method", func() {
		It("Should call buildResponse() function and returns responseData", func() {
			responseData := validationRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("question.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildResponseContent()))
			Expect(responseData.TemplateFormMethod).To(Equal("POST"))
			Expect(responseData.TemplateFormAction).To(Equal("/rule/networkPolicyRule"))
		})

		It("Should call nextRule Execute() method if the destinationAddress is a FQDN", func() {
			inputData.DestinationAddress = "sky.slack.com"
			responseData := validationRule.Execute(*inputData)
			Expect(responseData.HTTPStatus).To(Equal(200))
			Expect(responseData.TemplateName).To(Equal("response.tmpl"))
			Expect(responseData.TemplateContent).To(Equal(buildErrorResponseContent()))
		})
	})

	It("Should return a ResponseData with a message related to destinationAddress", func() {
		responseData := buildResponse(inputData.DestinationAddress)
		Expect(responseData.HTTPStatus).To(Equal(200))
		Expect(responseData.TemplateName).To(Equal("question.tmpl"))
		Expect(responseData.TemplateContent).To(Equal(buildResponseContent()))
		Expect(responseData.TemplateFormMethod).To(Equal("POST"))
		Expect(responseData.TemplateFormAction).To(Equal("/rule/networkPolicyRule"))
	})
})

func buildValidationRule() i.Rule {
	validationRule := new(ValidationRule)
	validationRule.SetName(c.VALIDATION_RULE)
	validationRule.SetNextRule(new(NetworkPolicyRule))
	return validationRule
}

func buildResponseContent() string {
	return fmt.Sprintf("Are you sure that your destination (%v) is an IP address and not a hostname? "+
		"The network filtering logic works based on how exactly "+
		"your applicaton reaches out to an external destination. If your "+
		"destination is configured as a raw IP, then you can continue!!", "192.168.1.0")
}

func buildErrorResponseContent() string {
	return "We apologize for the inconvenience, as we're currently encountering some technical issues. " +
		"Please get in touch with #core-support channel for further assistance."
}

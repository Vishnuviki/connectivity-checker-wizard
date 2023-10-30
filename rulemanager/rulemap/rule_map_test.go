package rulemap

import (
	"testing"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/rulemanager/rules"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRuleMap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Rule Map Suite")
}

var _ = Describe("Testing Rule Map functions", func() {
	var (
		ruleMap *RuleMap
	)

	BeforeEach(func() {
		ruleMap = GetInstance()
		ruleMap.m[c.VALIDATION_RULE] = new(rules.ValidationRule)
	})

	It("Should return same RuleMap instance", func() {
		instance := GetInstance()
		Expect(ruleMap).To(Equal(instance))
	})

	Context("Test GetRuleByName method", func() {
		It("Should return a existent rule", func() {
			rule, _ := ruleMap.GetRuleByName("validationRule")
			Expect(rule).NotTo(BeNil())
		})

		It("Should return a existent rule", func() {
			rule, _ := ruleMap.GetRuleByName("invalidRule")
			Expect(rule).To(BeNil())
		})
	})

	It("Shoud add rule into RuleMap", func() {
		ruleMap.AddRule(c.NETWORK_POLICY_RULE, new(rules.NetworkPolicyRule))
		rule, _ := ruleMap.GetRuleByName("networkPolicyRule")
		Expect(rule).NotTo(BeNil())
	})
})

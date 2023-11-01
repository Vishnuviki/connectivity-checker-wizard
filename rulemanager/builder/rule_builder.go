package builder

import (
	"log"

	"conectivity-checker-wizard/cilium"
	c "conectivity-checker-wizard/constants"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/rulemanager/rules"
)

var ruleMap = rulemap.GetInstance()

// This function build the rules based on the dependency order
func BuildRuleMap(policyChecker cilium.PolicyChecker) {
	buildDispatchIPRule()
	buildDNSLookUPRule()
	buildNetworkPolicyRule(policyChecker)
	buildValidationRule()
}

func getRuleByName(ruleName string) i.Rule {
	var rule i.Rule = nil
	if v, ok := ruleMap.GetRuleByName(ruleName); ok {
		rule = v
	} else {
		log.Printf("%s, is missing in the rule map", ruleName)
	}
	return rule
}

func buildDispatchIPRule() {
	rule := new(rules.DispatchIPRule)
	rule.SetName(c.DISPATCH_IP_RULE)
	rule.SetNextRule(nil)
	ruleMap.AddRule(c.DISPATCH_IP_RULE, rule)
}

// dnsLookUPRule has a dependency with dispatchIPRule
func buildDNSLookUPRule() {
	rule := new(rules.DNSLookUPRule)
	rule.SetName(c.DNS_LOOK_UP_RULE)
	rule.SetNextRule(getRuleByName(c.DISPATCH_IP_RULE))
	ruleMap.AddRule(c.DNS_LOOK_UP_RULE, rule)
}

func buildNetworkPolicyRule(policyChecker cilium.PolicyChecker) {
	rule := new(rules.NetworkPolicyRule)
	rule.SetName(c.NETWORK_POLICY_RULE)
	rule.SetNextRule(nil)
	rule.SetPolicyChecker(policyChecker)
	ruleMap.AddRule(c.NETWORK_POLICY_RULE, rule)
}

// validationRule has a dependency with networkPolicyRule
func buildValidationRule() {
	rule := new(rules.ValidationRule)
	rule.SetName(c.VALIDATION_RULE)
	rule.SetNextRule(getRuleByName(c.NETWORK_POLICY_RULE))
	ruleMap.AddRule(c.VALIDATION_RULE, rule)
}

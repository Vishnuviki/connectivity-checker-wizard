package builder

import (
	"log"

	c "conectivity-checker-wizard/constants"
	i "conectivity-checker-wizard/rulemanager/interfaces"
	"conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/rulemanager/rules"
	"conectivity-checker-wizard/services/cilium"
)

// This function build the rules based on the dependency order
func BuildRules(ruleMap *rulemap.RuleMap) {
	buildDispatchIPRule(ruleMap)
	buildDNSLookUPRule(ruleMap)
	buildNetworkPolicyRule(ruleMap)
	buildValidationRule(ruleMap)
}

func getRuleByName(ruleMap *rulemap.RuleMap, ruleName string) i.Rule {
	var rule i.Rule = nil
	if v, ok := ruleMap.GetRuleByName(ruleName); ok {
		rule = v
	} else {
		log.Printf("%s, is missing in the rule map", ruleName)
	}
	return rule
}

func buildDispatchIPRule(ruleMap *rulemap.RuleMap) {
	rule := new(rules.DispatchIPRule)
	rule.SetName(c.DISPATCH_IP_RULE)
	rule.SetNextRule(nil)
	ruleMap.AddRule(c.DISPATCH_IP_RULE, rule)
}

// dnsLookUPRule has a dependency with dispatchIPRule
func buildDNSLookUPRule(ruleMap *rulemap.RuleMap) {
	rule := new(rules.DNSLookUPRule)
	rule.SetName(c.DNS_LOOK_UP_RULE)
	rule.SetNextRule(getRuleByName(ruleMap, c.DISPATCH_IP_RULE))
	ruleMap.AddRule(c.DNS_LOOK_UP_RULE, rule)
}

func buildNetworkPolicyRule(ruleMap *rulemap.RuleMap) {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := rules.NewNetworkPolicyRule(c.NETWORK_POLICY_RULE, nil, cpc)
	rule.SetName(c.NETWORK_POLICY_RULE)
	rule.SetNextRule(nil)
	ruleMap.AddRule(c.NETWORK_POLICY_RULE, rule)
}

// validationRule has a dependency with networkPolicyRule
func buildValidationRule(ruleMap *rulemap.RuleMap) {
	rule := new(rules.ValidationRule)
	rule.SetName(c.VALIDATION_RULE)
	rule.SetNextRule(getRuleByName(ruleMap, c.NETWORK_POLICY_RULE))
	ruleMap.AddRule(c.VALIDATION_RULE, rule)
}

package builder

import (
	"log"

	c "conectivity-checker-wizard/constants"
	"conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/rulemanager/rules"
	"conectivity-checker-wizard/services/cilium"
)

func BuildRules(ruleMap *rulemap.RuleMap) {
	buildDispatchIPRule(ruleMap)
	buildDNSLookUPRule(ruleMap)
	buildNetworkPolicyRule(ruleMap)
	buildValidationRule(ruleMap)
}

func buildDispatchIPRule(ruleMap *rulemap.RuleMap) {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := rules.NewNetworkPolicyRule(c.NETWORK_POLICY_RULE, nil, cpc)
	ruleMap.AddRule(c.DISPATCH_IP_RULE, rule)
}

// dnsLookUPRule has a dependency with dispatchIPRule
func buildDNSLookUPRule(ruleMap *rulemap.RuleMap) {
	rule := new(rules.DNSLookUPRule)
	rule.SetName(c.DNS_LOOK_UP_RULE)
	if v, ok := ruleMap.GetRuleByName(c.DISPATCH_IP_RULE); ok {
		rule.SetNextRule(v)
	} else {
		log.Printf("%s, is missing in the rule map\n", c.DISPATCH_IP_RULE)
	}
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
	if v, ok := ruleMap.GetRuleByName(c.NETWORK_POLICY_RULE); ok {
		rule.SetNextRule(v)
	} else {
		log.Printf("%s, is missing in the rule map\n", c.NETWORK_POLICY_RULE)
	}
	ruleMap.AddRule(c.VALIDATION_RULE, rule)
}

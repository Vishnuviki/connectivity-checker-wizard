package builder

import (
	c "conectivity-checker-wizard/constants"
	rm "conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/services/cilium"
	r "conectivity-checker-wizard/services/rules"
)

var ruleMap = rm.GetInstance()

func BuildRules() {
	ruleMap.Map[c.DISPATCH_IP_RULE] = buildDispatchIPRule()
	ruleMap.Map[c.DNS_LOOK_UP_RULE] = buildDNSLookUPRule()
	ruleMap.Map[c.NETWORK_POLICY_RULE] = buildNetworkPolicyRule()
	ruleMap.Map[c.VALIDATION_RULE] = buildValidationRule()
}

func buildDispatchIPRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(c.NETWORK_POLICY_RULE, nil, cpc)
	return rule
}

func buildDNSLookUPRule() *r.DNSLookUPRule {
	rule := new(r.DNSLookUPRule)
	rule.SetName(c.DNS_LOOK_UP_RULE)
	if v, ok := ruleMap.GetRuleByName(c.DISPATCH_IP_RULE); ok {
		rule.SetNextRule(v)
	}
	return rule
}

func buildValidationRule() *r.ValidationRule {
	rule := new(r.ValidationRule)
	rule.SetName(c.VALIDATION_RULE)
	if v, ok := ruleMap.GetRuleByName(c.NETWORK_POLICY_RULE); ok {
		rule.SetNextRule(v)
	}
	return rule
}

func buildNetworkPolicyRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(c.NETWORK_POLICY_RULE, nil, cpc)
	rule.SetName(c.NETWORK_POLICY_RULE)
	rule.SetNextRule(nil)
	return rule
}

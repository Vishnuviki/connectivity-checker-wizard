package builder

import (
	rm "conectivity-checker-wizard/rulemanager/rulemap"
	"conectivity-checker-wizard/services/cilium"
	r "conectivity-checker-wizard/services/rules"
)

var ruleMap = rm.GetInstance()

func BuildRules() {
	ruleMap.Map[r.DISPATCH_IP_RULE] = buildDispatchIPRule()
	ruleMap.Map[r.DNS_LOOK_UP_RULE] = buildDNSLookUPRule()
	ruleMap.Map[r.NETWORK_POLICY_RULE] = buildNetworkPolicyRule()
	ruleMap.Map[r.VALIDATION_RULE] = buildValidationRule()
}

func buildDispatchIPRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(r.NETWORK_POLICY_RULE, nil, cpc)
	return rule
}

func buildDNSLookUPRule() *r.DNSLookUPRule {
	rule := new(r.DNSLookUPRule)
	rule.SetName(r.DNS_LOOK_UP_RULE)
	if v, ok := ruleMap.GetRuleByName(r.DISPATCH_IP_RULE); ok {
		rule.SetNextRule(v)
	}
	return rule
}

func buildValidationRule() *r.ValidationRule {
	rule := new(r.ValidationRule)
	rule.SetName(r.VALIDATION_RULE)
	if v, ok := ruleMap.GetRuleByName(r.NETWORK_POLICY_RULE); ok {
		rule.SetNextRule(v)
	}
	return rule
}

func buildNetworkPolicyRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(r.NETWORK_POLICY_RULE, nil, cpc)
	rule.SetName(r.NETWORK_POLICY_RULE)
	rule.SetNextRule(nil)
	return rule
}

package services

import (
	"conectivity-checker-wizard/models"
	r "conectivity-checker-wizard/rules"
	"conectivity-checker-wizard/cilium"

	"github.com/gin-gonic/gin"
)

var ruleMap = r.GetInstance()

func CreateRuleMap() {
	ruleMap.Map[r.DISPATCH_IP_RULE] = buildDispatchIPRule()
	ruleMap.Map[r.DNS_LOOK_UP_RULE] = buildDNSLookUPRule()
	ruleMap.Map[r.NETWORK_POLICY_RULE] = buildNetworkPolicyRule()
	ruleMap.Map[r.VALIDATION_RULE] = buildValidationRule()
}

func HandleRules(c *gin.Context, ruleName string) models.ResponseData {
	rule := ruleMap.GetRuleByName(ruleName)
	// execute rule
	return rule.Execute(c)
}

func buildDispatchIPRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(r.NETWORK_POLICY_RULE, nil, cpc)
	return rule
}

func buildDNSLookUPRule() *r.DNSLookUPRule {
	rule := new(r.DNSLookUPRule)
	rule.SetName(r.DNS_LOOK_UP_RULE)
	rule.SetNextRule(ruleMap.GetRuleByName(r.DISPATCH_IP_RULE))
	return rule
}

func buildNetworkPolicyRule() *r.NetworkPolicyRule {
	cpc := cilium.InClusterCiliumPolicyChecker{}
	rule := r.NewNetworkPolicyRule(r.NETWORK_POLICY_RULE, nil, cpc)
	rule.SetName(r.NETWORK_POLICY_RULE)
	rule.SetNextRule(nil)
	return rule
}

func buildValidationRule() *r.ValidationRule {
	rule := new(r.ValidationRule)
	rule.SetName(r.VALIDATION_RULE)
	rule.SetNextRule(nil)
	return rule
}

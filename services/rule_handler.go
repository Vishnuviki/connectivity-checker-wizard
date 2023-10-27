package services

import (
	"net/http"

	"conectivity-checker-wizard/cilium"
	"conectivity-checker-wizard/models"
	r "conectivity-checker-wizard/rules"

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
	if rule, ok := ruleMap.GetRuleByName(ruleName); ok {
		// execute rule
		return rule.Execute(c)
	} else {
		return r.BuildResponseData(http.StatusNotFound, "Page Not Found!!", "page-not-found.tmpl")
	}
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

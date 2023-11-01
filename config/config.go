package config

import (
	"os"

	"conectivity-checker-wizard/cilium"
	"conectivity-checker-wizard/rulemanager/builder"
)

// Congiure the system by building rules with the provided dependencies.
func Configure() {
	builder.BuildRuleMap(getCiliumPolicyChecker())
}

func getCiliumPolicyChecker() cilium.PolicyChecker {
	var policyChecker cilium.PolicyChecker
	if os.Getenv("LOCAL_DEV") == "true" {
		// use stubbed policy checker
		stub := &cilium.StubbedCiliumNetworkPolicyGetter{}
		policyChecker = cilium.NewInClusterCiliumPolicyChecker(stub)
	} else {
		policyChecker = cilium.NewInClusterCiliumPolicyChecker()
	}
	return policyChecker
}

package rules

import (
	"conectivity-checker-wizard/cilium"
	"fmt"
	"regexp"
)

func checkFQDNEgreesRules(namespace string, fqdn string) bool {
	policies, err := cilium.GetCiliumNetworkPolicies(namespace)
	if err != nil {
		return false
	}
	fqdnRegex := regexp.MustCompile(fqdn)
	for _, policy := range policies.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each FQDN in the egressRule
			for _, fqdn := range egressRule.ToFQDNs {
				val := fqdn.String()
				fmt.Println("FQDN:", val)
				if fqdnRegex.MatchString(val) {
					return true
				}
			}
		}
	}
	return false
}

func checkIPAddressEgressRules(namespace string, ip string) bool {
	policies, err := cilium.GetCiliumNetworkPolicies(namespace)
	if err != nil {
		return false
	}
	ipRegex := regexp.MustCompile(ip)
	for _, policy := range policies.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each CIDR in the egressRule
			for _, cidr := range egressRule.ToCIDR {
				val := string(cidr)
				fmt.Println("IP Address:", val)
				if ipRegex.MatchString(val) {
					return true
				}
			}
		}
	}
	return false
}

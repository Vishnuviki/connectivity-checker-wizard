package cilium

import (
	"fmt"
	"net"
	"strings"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

type CiliumPolicyChecker struct{}

func (c *CiliumPolicyChecker) checkIPAllowedByPolicyInNamespace(ip string, policyList *v2.CiliumNetworkPolicyList) bool {
	ipAddr := net.ParseIP(ip) // Parse the IP address string
	for _, policy := range policyList.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each CIDR in the egressRule
			for _, cidr := range egressRule.ToCIDR {
				_, ipNet, err := net.ParseCIDR(string(cidr))
				if err != nil {
					return false
				}
				if ipNet.Contains(ipAddr) {
					return true
				}
			}
		}
	}
	return false
}

func (c *CiliumPolicyChecker) checkFQDNAllowedByPolicyInNamespace(targetFQDN string, policyList *v2.CiliumNetworkPolicyList) bool {
	// Iterate through each network policy
	for _, policy := range policyList.Items {
		// Iterate through each egress rule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Iterate through each FQDN in the egress rule
			for _, fqdn := range egressRule.ToFQDNs {
				policyFQDN := fqdn.MatchPattern
				fmt.Println("FQDN FROM POLICY:", policyFQDN)
				if isFQDNMatchingPolicy(targetFQDN, policyFQDN) {
					return true

				}
			}
		}
	}
	return false
}

func isFQDNMatchingPolicy(targetFQDN, policyFQDN string) bool {
	targetFQDN = strings.ToLower(targetFQDN)
	policyFQDN = strings.ToLower(policyFQDN)
	// Split both FQDNs into parts
	targetParts := strings.Split(targetFQDN, ".")
	policyParts := strings.Split(policyFQDN, ".")
	// Check if the number of parts is the same or policyFQDN is a wildcard
	if len(targetParts) == len(policyParts) || strings.HasPrefix(policyFQDN, "*.") {
		return recursiveFQDNMatch(targetParts, policyParts)
	}
	return false
}

func recursiveFQDNMatch(targetParts, policyParts []string) bool {
	if len(policyParts) == 0 {
		return true
	}
	if policyParts[0] == "*" {
		// Handle wildcard, move to the next part in both FQDNs
		return recursiveFQDNMatch(targetParts[1:], policyParts[1:]) ||
			recursiveFQDNMatch(targetParts, policyParts[1:])
	}
	if policyParts[0] == targetParts[0] {
		// Match the current part, move to the next part
		return recursiveFQDNMatch(targetParts[1:], policyParts[1:])
	}
	return false
}

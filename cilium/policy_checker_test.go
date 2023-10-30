package cilium

import "testing"

func TestFQDNNameMatch(t *testing.T) {
	stub := &StubbedCiliumNetworkPolicyGetter{
		fqdnNames: []string{"host.domain.com"},
	}

	pc := NewInClusterCiliumPolicyChecker(stub)

	allowed, _ := pc.CheckFQDNAllowedByPolicyInNamespace("host.domain.com", "namespace")
	if !allowed {
		t.Error("FQDN should be allowed")
	}

	allowed, _ = pc.CheckFQDNAllowedByPolicyInNamespace("not-allowed.domain.com", "namespace")
	if allowed {
		t.Error("FQDN should not be allowed")
	}
}

func TestFQDNPatternMatch(t *testing.T) {
	stub := &StubbedCiliumNetworkPolicyGetter{
		fqdnPatterns: []string{"*.domain.com"},
	}

	pc := NewInClusterCiliumPolicyChecker(stub)

	allowed, _ := pc.CheckFQDNAllowedByPolicyInNamespace("host.domain.com", "namespace")
	if !allowed {
		t.Error("FQDN should be allowed")
	}

	allowed, _ = pc.CheckFQDNAllowedByPolicyInNamespace("host.potato.com", "namespace")
	if allowed {
		t.Error("FQDN should not be allowed")
	}
}
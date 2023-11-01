package cilium

import "testing"

func TestFQDNNameMatch(t *testing.T) {
	testCases := []struct {
		fqdnName    string
		fqdnToMatch string
		shouldMatch bool
	}{
		{fqdnName: "host.domain.com", fqdnToMatch: "host.domain.com", shouldMatch: true},
		{fqdnName: "host.domain.com", fqdnToMatch: "another.domain.com", shouldMatch: false},
	}

	for _, tc := range testCases {
		stub := &StubbedCiliumNetworkPolicyGetter{
			fqdnNames: []string{tc.fqdnName},
		}
		pc := NewInClusterCiliumPolicyChecker(stub)

		match, err := pc.CheckFQDNAllowedByPolicyInNamespace(tc.fqdnToMatch, "namespace")
		if err != nil {
			t.Errorf("Error checking FQDN %s, name: %s, error: %s", tc.fqdnToMatch, tc.fqdnName, err.Error())
		}
		if tc.shouldMatch && !match {
			t.Errorf("FQDN %s should match name %s, but it did not", tc.fqdnToMatch, tc.fqdnName)
		}
		if !tc.shouldMatch && match {
			t.Errorf("FQDN %s should not match name %s, but it did", tc.fqdnToMatch, tc.fqdnName)
		}
	}
}

// Robbed from: https://github.com/cilium/cilium/blob/5a0b88d1e0e4609c6f192a1b6aeadb46e2f48211/pkg/policy/api/fqdn.go#L62
//
// MatchPattern allows using wildcards to match DNS names. All wildcards are
// case insensitive. The wildcards are:
// - "*" matches 0 or more DNS valid characters, and may occur anywhere in
// the pattern. As a special case a "*" as the leftmost character, without a
// following "." matches all subdomains as well as the name to the right.
// A trailing "." is automatically added when missing.
//
// Examples:
// `*.cilium.io` matches subomains of cilium at that level
//
//	www.cilium.io and blog.cilium.io match, cilium.io and google.com do not
//
// `*cilium.io` matches cilium.io and all subdomains ends with "cilium.io"
//
//	except those containing "." separator, subcilium.io and sub-cilium.io match,
//	www.cilium.io and blog.cilium.io does not
//
// sub*.cilium.io matches subdomains of cilium where the subdomain component
// begins with "sub"
//
//	sub.cilium.io and subdomain.cilium.io match, www.cilium.io,
//	blog.cilium.io, cilium.io and google.com do not
func TestFQDNPatternMatch(t *testing.T) {
	testCases := []struct {
		fqdnPattern string
		fqdnToMatch string
		shouldMatch bool
	}{
		{fqdnPattern: "*", fqdnToMatch: "www.cilium.io", shouldMatch: true},

		{fqdnPattern: "*.cilium.io", fqdnToMatch: "www.cilium.io", shouldMatch: true},
		{fqdnPattern: "*.cilium.io", fqdnToMatch: "blog.cilium.io", shouldMatch: true},
		{fqdnPattern: "*.cilium.io", fqdnToMatch: "cilium.io", shouldMatch: false},
		{fqdnPattern: "*.cilium.io", fqdnToMatch: "abc.def.cilium.io", shouldMatch: false},

		{fqdnPattern: "abc.*.cilium.io", fqdnToMatch: "abc.def.cilium.io", shouldMatch: true},
		{fqdnPattern: "abc.def.cilium.io", fqdnToMatch: "abc.ghi.cilium.io", shouldMatch: false},

		{fqdnPattern: "*cilium.io", fqdnToMatch: "cilium.io", shouldMatch: true},
		{fqdnPattern: "*cilium.io", fqdnToMatch: "subcilium.io", shouldMatch: true},
		{fqdnPattern: "*cilium.io", fqdnToMatch: "www.cilium.io", shouldMatch: false},
		{fqdnPattern: "*cilium.io", fqdnToMatch: "blog.cilium.io", shouldMatch: false},

		{fqdnPattern: "sub*.cilium.io", fqdnToMatch: "sub.cilium.io", shouldMatch: true},
		{fqdnPattern: "sub*.cilium.io", fqdnToMatch: "subdomain.cilium.io", shouldMatch: true},
		{fqdnPattern: "sub*.cilium.io", fqdnToMatch: "blog.cilium.io", shouldMatch: false},
		{fqdnPattern: "sub*.cilium.io", fqdnToMatch: "www.cilium.io", shouldMatch: false},

		{fqdnPattern: "www.cilium.io", fqdnToMatch: "www.cilium.io", shouldMatch: true},
		{fqdnPattern: "www.cilium.io", fqdnToMatch: "zzz.cilium.io", shouldMatch: false},

		{fqdnPattern: "no-dots", fqdnToMatch: "no-dots", shouldMatch: false},
	}

	for _, tc := range testCases {
		stub := &StubbedCiliumNetworkPolicyGetter{
			fqdnPatterns: []string{tc.fqdnPattern},
		}
		pc := NewInClusterCiliumPolicyChecker(stub)

		match, err := pc.CheckFQDNAllowedByPolicyInNamespace(tc.fqdnToMatch, "namespace")
		if err != nil {
			t.Errorf("Error checking FQDN %s, pattern: \"%s\", error: %s", tc.fqdnToMatch, tc.fqdnPattern, err.Error())
		}
		if tc.shouldMatch && !match {
			t.Errorf("FQDN %s should match pattern \"%s\", but it did not", tc.fqdnToMatch, tc.fqdnPattern)
		}
		if !tc.shouldMatch && match {
			t.Errorf("FQDN %s should not match pattern \"%s\", but it did", tc.fqdnToMatch, tc.fqdnPattern)
		}
	}
}

func TestCIDRMatch(t *testing.T) {
	testCases := []struct {
		cidr        string
		ipToMatch   string
		shouldMatch bool
	}{
		{cidr: "10.0.0.0/8", ipToMatch: "10.0.0.10", shouldMatch: true},
		{cidr: "10.0.0.0/8", ipToMatch: "192.168.0.10", shouldMatch: false},
	}

	for _, tc := range testCases {
		stub := &StubbedCiliumNetworkPolicyGetter{
			cidrs: []string{tc.cidr},
		}
		pc := NewInClusterCiliumPolicyChecker(stub)

		match, err := pc.CheckIPAllowedByPolicyInNamespace(tc.ipToMatch, "namespace")
		if err != nil {
			t.Errorf("Error checking IP %s, cidr range: %s, error: %s", tc.ipToMatch, tc.cidr, err.Error())
		}
		if tc.shouldMatch && !match {
			t.Errorf("IP %s should match cidr range %s, but it did not", tc.ipToMatch, tc.cidr)
		}
		if !tc.shouldMatch && match {
			t.Errorf("IP %s should not match cidr range %s, but it did", tc.ipToMatch, tc.cidr)
		}
	}
}

func TestCIDRSetMatch(t *testing.T) {
	testCases := []struct {
		cidrSet     cidrSet
		ipToMatch   string
		shouldMatch bool
	}{
		{
			cidrSet:     cidrSet{cidr: "0.0.0.0/0", exceptCidrs: []string{"172.21.20.0/22"}},
			ipToMatch:   "10.0.0.10",
			shouldMatch: true,
		},
		{
			cidrSet:     cidrSet{cidr: "0.0.0.0/0", exceptCidrs: []string{"172.21.20.0/22"}},
			ipToMatch:   "172.21.20.10",
			shouldMatch: false,
		},
		{
			cidrSet:     cidrSet{cidr: "10.0.0.0/8", exceptCidrs: []string{}},
			ipToMatch:   "172.21.20.10",
			shouldMatch: false,
		},
		{
			cidrSet:     cidrSet{cidr: "172.31.0.0/16", exceptCidrs: []string{"172.21.20.0/22"}},
			ipToMatch:   "172.21.20.10",
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		stub := &StubbedCiliumNetworkPolicyGetter{
			cidrSets: []cidrSet{tc.cidrSet},
		}
		pc := NewInClusterCiliumPolicyChecker(stub)

		match, err := pc.CheckIPAllowedByPolicyInNamespace(tc.ipToMatch, "namespace")
		if err != nil {
			t.Errorf("Error checking IP %s, cidrSet: %s, error: %s", tc.ipToMatch, tc.cidrSet, err.Error())
		}
		if tc.shouldMatch && !match {
			t.Errorf("IP %s should match cidrSet %s, but it did not", tc.ipToMatch, tc.cidrSet)
		}
		if !tc.shouldMatch && match {
			t.Errorf("IP %s should not match cidrSet %s, but it did", tc.ipToMatch, tc.cidrSet)
		}
	}
}

package cilium

import (
	"context"
	"log"
	"regexp"
	"strings"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	cilium_clientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	cilium_v2 "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/typed/cilium.io/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

type PolicyChecker interface {
	CheckIPAllowedByPolicyInNamespace(ip string, namespace string) (bool, error)
	CheckFQDNAllowedByPolicyInNamespace(fqdn string, namespace string) (bool, error)
}

type InClusterCiliumPolicyChecker struct {
	getter cilium_v2.CiliumNetworkPoliciesGetter
}

func NewInClusterCiliumPolicyChecker(getter ...cilium_v2.CiliumNetworkPoliciesGetter) PolicyChecker {
	var g cilium_v2.CiliumNetworkPoliciesGetter

	// allow stubbed getter for testing
	if len(getter) > 0 {
		g = getter[0]
	}

	if g == nil {
		g = inClusterCiliumNetworkPoliciesGetter()
	}

	return &InClusterCiliumPolicyChecker{
		getter: g,
	}
}

func (c *InClusterCiliumPolicyChecker) CheckIPAllowedByPolicyInNamespace(ip string, namespace string) (bool, error) {
	policies, err := c.getCiliumNetworkPolicies(namespace)
	if err != nil {
		return false, err
	}

	ipRegex := regexp.MustCompile(ip)
	for _, policy := range policies.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each CIDR in the egressRule
			for _, cidr := range egressRule.ToCIDR {
				val := string(cidr)
				log.Println("IP Address:", val)
				if ipRegex.MatchString(val) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (c *InClusterCiliumPolicyChecker) CheckFQDNAllowedByPolicyInNamespace(fqdn string, namespace string) (bool, error) {
	policies, err := c.getCiliumNetworkPolicies(namespace)
	if err != nil {
		return false, err
	}

	for _, policy := range policies.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each FQDN in the egressRule
			for _, fqdnSelector := range egressRule.ToFQDNs {
				if isNameMatch(fqdn, fqdnSelector.MatchName) {
					return true, nil
				}
				if isPatternMatch(fqdn, fqdnSelector.MatchPattern) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func isNameMatch(fqdn string, matchName string) bool {
	return strings.EqualFold(fqdn, matchName)
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
func isPatternMatch(fqdn string, matchPattern string) bool {
	if matchPattern == "" {
		return false
	}

	patternTokens := strings.Split(matchPattern, ".")
	fqdnTokens := strings.Split(fqdn, ".")

	// TODO: do we allow "*" name pattern? Feels like this should not be allowed,
	// TODO: otherwise, what is the point if you can just say "*" and everything matches?
	// TODO: double check!

	// no "." in pattern or fqdn
	if len(patternTokens) == 1 || len(fqdnTokens) == 1 {
		return false
	}

	if len(patternTokens) != len(fqdnTokens) {
		return false
	}

	// iterating from the end
	for i := len(patternTokens) - 1; i >= 0; i-- {
		patternToken := patternTokens[i]
		fqdnToken := fqdnTokens[i]

		// not the first token
		if i > 0 {
			// TODO: do we allow "*sub" or "sub*" matches here?
			if patternToken == "*" {
				continue
			}
			if patternToken != fqdnToken {
				return false
			}
		}

		// first token
		if i == 0 {
			if patternToken == "*" {
				return true
			}
			// "*sub" case, match *subdomain and subdomain
			if strings.HasPrefix(patternToken, "*") {
				return strings.HasSuffix(fqdnToken, patternToken[1:]) ||
					fqdnToken == patternToken[1:]
			}
			//sub* case
			if strings.HasSuffix(patternToken, "*") {
				return strings.HasPrefix(fqdnToken, patternToken[:len(patternToken)-1])
			}
			if patternToken != fqdnToken {
				return false
			}
		}
	}

	return true
}

func (c *InClusterCiliumPolicyChecker) getCiliumNetworkPolicies(namespace string) (*v2.CiliumNetworkPolicyList, error) {
	policies, err := c.getter.CiliumNetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error reading cilium network policies", err)
		return nil, err
	}

	return policies, nil
}

func clientset() (*cilium_clientset.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Error creating in-cluster config", err)
		return nil, err
	}

	// creates the clientset
	clientset, err := cilium_clientset.NewForConfig(config)
	if err != nil {
		log.Println("Error creating clientset", err)
		return nil, err
	}

	return clientset, nil
}

func inClusterCiliumNetworkPoliciesGetter() cilium_v2.CiliumNetworkPoliciesGetter {
	clientset, err := clientset()
	if err != nil {
		log.Fatalln("Error creating in-cluster clientset, are you in a cluster? Exiting!", err)
	}
	return clientset.CiliumV2()
}

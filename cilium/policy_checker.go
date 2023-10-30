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
	gonet "github.com/THREATINT/go-net"
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
				if strings.EqualFold(fqdnSelector.MatchName, fqdn) {
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

// MatchPattern allows using wildcards to match DNS names. All wildcards are
// case insensitive. The wildcards are:
// - "*" matches 0 or more DNS valid characters, and may occur anywhere in
// the pattern. As a special case a "*" as the leftmost character, without a
// following "." matches all subdomains as well as the name to the right.
// A trailing "." is automatically added when missing.
//
// Examples:
// `*.cilium.io` matches subomains of cilium at that level
//   www.cilium.io and blog.cilium.io match, cilium.io and google.com do not
// `*cilium.io` matches cilium.io and all subdomains ends with "cilium.io"
//   except those containing "." separator, subcilium.io and sub-cilium.io match,
//   www.cilium.io and blog.cilium.io does not
// sub*.cilium.io matches subdomains of cilium where the subdomain component
// begins with "sub"
//   sub.cilium.io and subdomain.cilium.io match, www.cilium.io,
//   blog.cilium.io, cilium.io and google.com do not
//
func isPatternMatch(fqdn string, matchPattern string) bool {
	if matchPattern == "" {
		return false
	}

	if gonet.DomainFromFqdn(fqdn) == gonet.DomainFromFqdn(matchPattern) {
		return true
	}

	return false
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

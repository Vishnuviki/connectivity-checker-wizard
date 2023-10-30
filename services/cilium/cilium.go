package cilium

import (
	"context"
	"log"
	"regexp"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	cilium_clientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

type CiliumPolicyChecker interface {
	CheckIPAllowedByPolicyInNamespace(ip string, namespace string) (bool, error)
	CheckFQDNAllowedByPolicyInNamespace(fqdn string, namespace string) (bool, error)
}

type InClusterCiliumPolicyChecker struct{}

func (c InClusterCiliumPolicyChecker) CheckIPAllowedByPolicyInNamespace(ip string, namespace string) (bool, error) {
	policies, err := GetCiliumNetworkPolicies(namespace)
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

func (c InClusterCiliumPolicyChecker) CheckFQDNAllowedByPolicyInNamespace(fqdn string, namespace string) (bool, error) {
	policies, err := GetCiliumNetworkPolicies(namespace)
	if err != nil {
		return false, err
	}
	fqdnRegex := regexp.MustCompile(fqdn)
	for _, policy := range policies.Items {
		// Check each egressRule in the policy
		for _, egressRule := range policy.Spec.Egress {
			// Check each FQDN in the egressRule
			for _, fqdn := range egressRule.ToFQDNs {
				val := fqdn.String()
				log.Println("FQDN:", val)
				if fqdnRegex.MatchString(val) {
					return true, nil
				}
			}
		}
	}
	return false, nil
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

// TODO: this should be private
func GetCiliumNetworkPolicies(namespace string) (*v2.CiliumNetworkPolicyList, error) {
	clientset, err := clientset()
	if err != nil {
		return nil, err
	}

	policies, err := clientset.CiliumV2().CiliumNetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error reading cilium network policies", err)
		return nil, err
	}

	return policies, nil
}

package cilium

import (
	"context"
	"log"

	"github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	clientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type CiliumClient struct {
	policyChecker *CiliumPolicyChecker
}

func NewCiliumClient() CiliumClient {
	return CiliumClient{
		policyChecker: new(CiliumPolicyChecker),
	}
}

func (c *CiliumClient) GetCiliumNetworkPolicies(namespace string) (*v2.CiliumNetworkPolicyList, error) {
	clientSet, err := getInClusterClientSet()
	if err != nil {
		log.Println("Error accessing InClusterClientSet", err)
		return nil, err
	}
	policies, err := clientSet.CiliumV2().CiliumNetworkPolicies(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Println("Error reading cilium network policies", err)
		return nil, err
	}
	return policies, nil
}

func (c *CiliumClient) CheckIPAllowedByPolicyInNamespace(ip string, policyList *v2.CiliumNetworkPolicyList) bool {
	return c.policyChecker.checkIPAllowedByPolicyInNamespace(ip, policyList)
}

func (c *CiliumClient) CheckFQDNAllowedByPolicyInNamespace(ip string, policyList *v2.CiliumNetworkPolicyList) bool {
	return c.policyChecker.checkFQDNAllowedByPolicyInNamespace(ip, policyList)
}

func getInClusterClientSet() (*clientset.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Error creating in-cluster config", err)
		return nil, err
	}
	// creates the clientset
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		log.Println("Error creating clientset", err)
		return nil, err
	}
	return clientset, nil
}

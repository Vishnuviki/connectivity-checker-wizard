package cilium

import (
	"context"
	"log"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	cilium_clientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

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

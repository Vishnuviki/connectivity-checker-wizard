package cilium

import (
	"context"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	cilium_v2 "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/typed/cilium.io/v2"
	"github.com/cilium/cilium/pkg/policy/api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
)

type StubbedCiliumNetworkPolicyGetter struct {
	fqdnNames []string
	fqdnPatterns []string
	cidrs []string
}

type StubbedCiliumNetworkPolicyInterface struct {
	stub *StubbedCiliumNetworkPolicyGetter
}

var _ cilium_v2.CiliumNetworkPoliciesGetter = (*StubbedCiliumNetworkPolicyGetter)(nil)

func (s *StubbedCiliumNetworkPolicyGetter) CiliumNetworkPolicies(namespace string) cilium_v2.CiliumNetworkPolicyInterface {
	return &StubbedCiliumNetworkPolicyInterface{stub: s}
}

var _ cilium_v2.CiliumNetworkPolicyInterface = (*StubbedCiliumNetworkPolicyInterface)(nil)

func (s *StubbedCiliumNetworkPolicyInterface) List(ctx context.Context, opts v1.ListOptions) (*v2.CiliumNetworkPolicyList, error) {
	cpl := &v2.CiliumNetworkPolicyList{
		Items: []v2.CiliumNetworkPolicy{
			{
				Spec: &api.Rule{
					Egress: []api.EgressRule{
						{
							ToFQDNs: []api.FQDNSelector{},
							EgressCommonRule: api.EgressCommonRule{},
						},
					},
				},
			},
		},
	}

	if len(s.stub.fqdnNames) > 0 || len(s.stub.fqdnPatterns) > 0 {
		cpl.Items[0].Spec.Egress[0].ToFQDNs = toFQDNs(s.stub.fqdnNames, s.stub.fqdnPatterns)
	}

	if len(s.stub.cidrs) > 0 {
		cpl.Items[0].Spec.Egress[0].ToCIDR = toCIDRs(s.stub.cidrs)
	}

	return cpl, nil
}

func toFQDNs(fqdnNames []string, fqdnPatterns []string) []api.FQDNSelector {
	var fqdnSelectors []api.FQDNSelector
	for _, name := range fqdnNames {
		fqdnSelectors = append(fqdnSelectors, api.FQDNSelector{MatchName: name})
	}
	for _, pattern := range fqdnPatterns {
		fqdnSelectors = append(fqdnSelectors, api.FQDNSelector{MatchPattern: pattern})
	}
	return fqdnSelectors
}

func toCIDRs(cidrs []string) []api.CIDR {
	var cidrSelectors []api.CIDR
	for _, c := range cidrs {
		cidrSelectors = append(cidrSelectors, api.CIDR(c))
	}
	return cidrSelectors
}

// example network policy:
// var defaultPolicyList *v2.CiliumNetworkPolicyList = &v2.CiliumNetworkPolicyList{
// 	Items: []v2.CiliumNetworkPolicy{
// 		{
// 			Spec: &api.Rule{
// 				Egress: []api.EgressRule{
// 					{
// 						ToFQDNs: []api.FQDNSelector{
// 							{
// 								MatchName: "www.bbc.co.uk",
// 								MatchPattern: "*.slack.com",
// 							},
// 						},
// 						EgressCommonRule: api.EgressCommonRule{
// 							ToCIDR: []api.CIDR{
// 								"0.0.0.0/0",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

// ignore these, only to satisfy the interface, we only care about the List method
func (s *StubbedCiliumNetworkPolicyInterface) Create(ctx context.Context, ciliumNetworkPolicy *v2.CiliumNetworkPolicy, opts v1.CreateOptions) (*v2.CiliumNetworkPolicy, error) {return nil, nil}
func (s *StubbedCiliumNetworkPolicyInterface) Update(ctx context.Context, ciliumNetworkPolicy *v2.CiliumNetworkPolicy, opts v1.UpdateOptions) (*v2.CiliumNetworkPolicy, error) {return nil, nil}
func (s *StubbedCiliumNetworkPolicyInterface) UpdateStatus(ctx context.Context, ciliumNetworkPolicy *v2.CiliumNetworkPolicy, opts v1.UpdateOptions) (*v2.CiliumNetworkPolicy, error) {return nil, nil}
func (s *StubbedCiliumNetworkPolicyInterface) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {return nil}
func (s *StubbedCiliumNetworkPolicyInterface) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {return nil}
func (s *StubbedCiliumNetworkPolicyInterface) Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.CiliumNetworkPolicy, error) {return nil, nil}
func (s *StubbedCiliumNetworkPolicyInterface) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {return nil, nil}
func (s *StubbedCiliumNetworkPolicyInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.CiliumNetworkPolicy, err error) {return nil, nil}

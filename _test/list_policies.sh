#!/usr/bin/env bash

set -euo pipefail

# Run the below command for a quick glance over all egress policies in a cluster.
# We ignore "toEndpoints" as these configure egress to endpoints within a cluster.
# That leaves us with: toCIDR, toCIDRSet, toFQDNs to deal with.
#
# kubectl get ciliumnetworkpolicies.cilium.io \
#     --all-namespaces \
#     -o=jsonpath='{range .items[*].spec.egress[*]}{@}{"\n"}{end}' \
#     | sort | uniq | less

list_toFQDNs() {
    kubectl get ciliumnetworkpolicies.cilium.io \
        --all-namespaces \
        -o=jsonpath='{range .items[*].spec.egress[*].toFQDNs[*]}{@}{"\n"}{end}' \
        | sort | uniq
}

list_toCIDR() {
    kubectl get ciliumnetworkpolicies.cilium.io \
        --all-namespaces \
        -o=jsonpath='{range .items[*].spec.egress[*].toCIDR[*]}{@}{"\n"}{end}' \
        | sort | uniq
}

list_toCIDRSet() {
    kubectl get ciliumnetworkpolicies.cilium.io \
        --all-namespaces \
        -o=jsonpath='{range .items[*].spec.egress[*].toCIDRSet[*]}{@}{"\n"}{end}' \
        | sort | uniq
}

echo "FQDNs *********************************************************"
list_toFQDNs

echo "CIDR *********************************************************"
list_toCIDR

echo "CIDRSet *********************************************************"
list_toCIDRSet

# Conectivity Wizard

## Try it locally

```
make run
```
Visit http://localhost:8080

## Try it from local k8s

Install kind and skaffold:

```
brew install kind skaffold
```

Create a kind cluster:

```
kind create cluster \
    --name kind \
    --image kindest/node:v1.23.17 \
    --config kind-config.yaml

# install cilium
helm repo add cilium https://helm.cilium.io/
docker pull quay.io/cilium/cilium:v1.14.3
kind load docker-image quay.io/cilium/cilium:v1.14.3

helm install cilium cilium/cilium --version 1.14.3 \
   --namespace kube-system \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes \
   --set operator.replicas=1

# check cilium status
cilium status --wait
```

Grant admin cluster role to default:default service acount:

```
# create cilium cluster role
kubectl apply -f _test/cilium-cluster-role.yaml

# create bindings
kubectl create clusterrolebinding default-admin \
 	--clusterrole=admin  \
 	--serviceaccount=default:default

kubectl create clusterrolebinding default-cilium \
 	--clusterrole=cilium  \
 	--serviceaccount=default:default
```

Create cilium network policies:

```
# network policy in default ns
kubectl apply -f _test/network-policy.yaml

# cluster wide policies allowing traffic internal to the cluster
kubectl apply -f _test/cluster-network-policy.yaml
```

Deploy the application to your kind cluster, make sure your current namespace is `default`:

```
skaffold run
```

Port-forward the application to localhost:8080:

```
kubectl port-forward service/web 8080:8080
```

Test the application:

```
curl localhost:8080
```

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
kind create cluster --name kind --image kindest/node:v1.23.17
```

Deploy the application to your local cluster in the default namespace:

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

# canivete-api k8s

## Deployment

Create the namespace

```
kubectl create ns canivete
```

Install Helm chart

```
helm install canivete-api-chart ./api-chart --namespace canivete
```

Delete Helm chart

```
helm uninstall canivete-api-chart --namespace canivete
```

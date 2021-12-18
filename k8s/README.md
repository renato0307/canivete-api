# canivete-api k8s

## Deployment

Create the namespace

```
kubectl create ns canivete
```

Create the deployment and service

```
kubectl apply -f canivete-rest-deployment.yaml
kubectl apply -f canivete-rest-service.yaml
```

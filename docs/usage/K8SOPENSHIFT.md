# K8s and OpenShift

## K8s

```bash
kubectl create ns yaam
```

## OpenShift

```bash
oc new-project yaam
```

## Deploy

```bash
kubectl create -f deployments/k8s-openshift/deploy.yml -n yaam
```

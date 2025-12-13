### API Gateway

```
User -> API Gateway -> Backend
```

```
helm upgrade -i prometheus oci://ghcr.io/prometheus-community/charts/kube-prometheus-stack -n monitoring --values values.yaml
```

```
minikube addons enable ingress
```

<p align="center">
  <img src="apigw.drawio.png" alt="API Gateway diagram" width="400"/>
</p>
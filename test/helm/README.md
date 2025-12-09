```
eval $(minikube docker-env)
eval $(minikube docker-env --unset)
helm upgrade -i gateway . -f  values.yaml
```
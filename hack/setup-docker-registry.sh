kubectl apply -f hack/kube-docker-registry.yaml
kubectl wait pod/$(kubectl get po -n kube-system | grep kube-registry-v0 | awk '{print $1;}') -n kube-system --for condition=Ready --timeout=180s
kubectl port-forward --namespace kube-system $(kubectl get po -n kube-system | grep kube-registry-v0 | awk '{print $1;}') 30500:5000

# metal lb
kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml
kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/example-layer2-config.yaml


# nginx infress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml


# rook.io
kubectl apply -f https://raw.githubusercontent.com/rook/rook/release-0.9/cluster/examples/kubernetes/ceph/operator.yaml
kubectl apply -f https://raw.githubusercontent.com/rook/rook/release-0.9/cluster/examples/kubernetes/ceph/cluster.yaml

# etcd operator
kubectl apply -f hack/dev/etcd-operator-rbac.yml
kubectl apply -f hack/dev/etcd-operator-deployment.yml

# wait for crd
NEXT_WAIT_TIME=0
until kubectl get EtcdCluster || [ $NEXT_WAIT_TIME -eq 60 ]; do
   sleep $(( NEXT_WAIT_TIME++ ))
done

# etcd cluster
kubectl apply -f hack/dev/coredns-etcd-cluster.yml

# external dns
kubectl apply -f hack/dev/extdns.yml

# core dns
kubectl apply -f hack/dev/coredns-extdns.yml

# virtlet lb
kubectl apply -f https://raw.githubusercontent.com/ivan4th/virtletlb/master/config/crds/virtletlb_v1alpha1_innerservice.yaml
kubectl apply -f https://raw.githubusercontent.com/lukaszo/cluster-api-provider-virtlet/master/hack/examples/outer-controller.yaml

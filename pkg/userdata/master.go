package userdata

const masterProvisionScript = `
- path: /usr/local/bin/provision.sh
  permissions: "0755"
  owner: root
  content: |
    #!/bin/bash
    set -u -e
    set -o pipefail
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
    apt-get update
    apt-get install -y docker.io kubelet kubeadm kubectl kubernetes-cni ceph-common python-pip
    sed -i 's/--cluster-dns=10\.96\.0\.10/--cluster-dns=10.97.0.10/' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
    systemctl daemon-reload

    # TODO generate token, use cidrs from cluster
    kubeadm init --token adcb82.4eae29627dc4c5a6 --pod-network-cidr=10.200.0.0/16 --service-cidr=10.97.0.0/16 --apiserver-cert-extra-sans=127.0.0.1,localhost

    # TODO: installing weve, customize it
    export KUBECONFIG=/etc/kubernetes/admin.conf
    export kubever=$(kubectl version | base64 | tr -d '\n')
    kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$kubever"
    while ! kubectl get pods -n kube-system -l k8s-app=kube-dns|grep ' 1/1'; do
      sleep 1
    done
    mkdir -p /root/.kube
    chmod 700 /root/.kube
    cp "${KUBECONFIG}" /root/.kube/config

    # configure ceph storage
    kubectl apply -f /root/ceph.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/serviceaccount.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/role.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/clusterrole.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/clusterrolebinding.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/rolebinding.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-incubator/external-storage/master/ceph/rbd/deploy/rbac/deployment.yaml

    # LB Controller 'inner' part
    kubectl apply -f https://raw.githubusercontent.com/lukaszo/cluster-api-provider-virtlet/master/hack/examples/inner-controller.yaml

    # FIXME: enable insecure port
    # IT HAS TO BE LAST ACTION
    sed -i "s/--insecure-port=0/--insecure-port=8080\\n    - --insecure-bind-address=0.0.0.0/" /etc/kubernetes/manifests/kube-apiserver.yaml
    echo "Master setup complete." >&2
`

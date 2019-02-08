package userdata

const computeProvisionScript = `
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

    kubeadm join --token adcb82.4eae29627dc4c5a6 --discovery-token-unsafe-skip-ca-verification api-server:6443
    echo "Node setup complete." >&2
`

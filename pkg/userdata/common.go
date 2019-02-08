package userdata

import (
	"fmt"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	providerv1 "sigs.k8s.io/cluster-api-provider-virtlet/pkg/apis/virtlet/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func getDockerService() string {
	return `
- path: /etc/systemd/system/docker.service.d/env.conf
  permissions: "0644"
  owner: root
  content: |
    [Service]
    Environment="DOCKER_OPTS=--storage-driver=overlay2"
`
}

func getKubernetesAPTKey() string {
	return `
- path: /etc/apt/sources.list.d/kubernetes.list
  permissions: "0644"
  owner: root
  content: |
    deb http://apt.kubernetes.io/ kubernetes-xenial main
`
}

func getUsers() string {
	return `
- name: root
  # VirtletSSHKeys only affects 'ubuntu' user for this image, but we want root access
  ssh-authorized-keys:
  - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCaJEcFDXEK2ZbX0ZLS1EIYFZRbDAcRfuVjpstSc0De8+sV1aiu+dePxdkuDRwqFtCyk6dEZkssjOkBXtri00MECLkir6FcH3kKOJtbJ6vy3uaJc9w1ERo+wyl6SkAh/+JTJkp7QRXj8oylW5E20LsbnA/dIwWzAF51PPwF7A7FtNg9DnwPqMkxFo1Th/buOMKbP5ZA1mmNNtmzbMpMfJATvVyiv3ccsSJKOiyQr6UG+j7sc/7jMVz5Xk34Vd0l8GwcB0334MchHckmqDB142h/NCWTr8oLakDNvkfC1YneAfAO41hDkUbxPtVBG5M/o7P4fxoqiHEX+ZLfRxDtHB53 me@localhost
`
}

func getCephFiles(clientset *kubernetes.Clientset, cluster *clusterv1.Cluster) string {
	providerConf, err := providerv1.ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		log.Printf("Couldn't get providerSpec from cluster: %s: %v", cluster.Name, err)
		return ""
	}
	adminKey := providerConf.CephAdminKey
	clientKey := providerConf.CephClientKey

	svcs, err := clientset.CoreV1().Services("rook-ceph").List(metav1.ListOptions{LabelSelector: "app=rook-ceph-mon"})
	if err != nil {
		log.Printf("Couldn not list monitors in 'rook-ceph' namespace: %v", err)
		return ""
	}
	monitors := []string{}
	for _, svc := range svcs.Items {
		monitors = append(monitors, svc.Spec.ClusterIP+":6790")
	}

	specs := `
- path: /root/ceph.yaml
  permissions: "0600"
  owner: root
  content: |
    apiVersion: v1
    kind: Secret
    metadata:
      name: ceph-admin-secret
      namespace: default
    type: "kubernetes.io/rbd"
    data:
      # ceph auth get-key client.admin | base64
      key: %s
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: ceph-secret
      namespace: default
    type: "kubernetes.io/rbd"
    data:
      # ceph auth add client.kube mon 'allow r' osd 'allow rwx pool=kube'
      # ceph auth get-key client.kube | base64
      key: %s
    ---
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
      name: rbd
      annotations:
        storageclass.kubernetes.io/is-default-class: "true"
    provisioner: ceph.com/rbd
    parameters:
      monitors: %s
      pool: %s
      adminId: admin
      adminSecretNamespace: default
      adminSecretName: ceph-admin-secret
      userId: kube
      userSecretNamespace: default
      userSecretName: ceph-secret
      imageFormat: "2"
      imageFeatures: layering
`
	return fmt.Sprintf(specs, adminKey, clientKey, strings.Join(monitors, ","), cluster.Name)
}

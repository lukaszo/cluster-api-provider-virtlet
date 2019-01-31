/*
Copyright 2019 Mirantis.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machine

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	providerv1 "sigs.k8s.io/cluster-api-provider-virtlet/pkg/apis/virtlet/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"
)

const (
	ProviderName = "virtlet"
)

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	machinesGetter client.MachinesGetter
	clientset      *kubernetes.Clientset
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	MachinesGetter client.MachinesGetter
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Actuator{
		machinesGetter: params.MachinesGetter,
		clientset:      clientset,
	}, nil
}

// Create creates a machine and is invoked by the Machine Controller
func (a *Actuator) Create(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	log.Printf("Creating machine %v for cluster %v.", machine.Name, cluster.Name)
	err := a.reconcilePod(ctx, cluster, machine)
	if err != nil {
		return fmt.Errorf("Could not create machine for cluster %s: %v", cluster.Name, err)
	}
	return nil
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	log.Printf("Deleting machine %v for cluster %v.", machine.Name, cluster.Name)
	_, err := a.clientset.CoreV1().Pods(cluster.Namespace).Get(machine.Name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("Could not get machine pod: %v", err)
	}
	err = a.clientset.CoreV1().Pods(cluster.Namespace).Delete(machine.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Could not delete machine pod: %v", err)
	}
	return nil
}

// Update updates a machine and is invoked by the Machine Controller
func (a *Actuator) Update(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	log.Printf("Updating machine %v for cluster %v.", machine.Name, cluster.Name)
	return nil
}

// Exists test for the existance of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	log.Printf("Checking if machine %v for cluster %v exists.", machine.Name, cluster.Name)
	_, err := a.clientset.CoreV1().Pods(cluster.Namespace).Get(machine.Name, metav1.GetOptions{})
	if err != nil {
		// TODO: check if error is different from "not found"
		return false, nil
	}
	return true, nil
}

// The Machine Actuator interface must implement GetIP and GetKubeConfig functions as a workaround for issues
// cluster-api#158 (https://github.com/kubernetes-sigs/cluster-api/issues/158) and cluster-api#160
// (https://github.com/kubernetes-sigs/cluster-api/issues/160).

// GetIP returns IP address of the machine in the cluster.
func (a *Actuator) GetIP(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	log.Printf("Getting IP of machine %v for cluster %v.", machine.Name, cluster.Name)
	pod, err := a.clientset.CoreV1().Pods(cluster.Namespace).Get(machine.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Getting IP for pod (%s) for cluster %v failed: %v.", machine.Name, cluster.Name, err)
		return "", fmt.Errorf("Could not get IP of the pod (%s) for cluster %s: %v", machine.Name, cluster.Name, err)
	}

	return pod.Status.PodIP, nil
}

// GetKubeConfig gets a kubeconfig from the master.
func (a *Actuator) GetKubeConfig(cluster *clusterv1.Cluster, master *clusterv1.Machine) (string, error) {
	log.Printf("Getting IP of machine %v for cluster %v.", master.Name, cluster.Name)
	return "", fmt.Errorf("TODO: Not yet implemented")
}

func (a *Actuator) reconcilePod(ctx context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	_, err := a.clientset.CoreV1().Pods(cluster.Namespace).Get(machine.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Creating pod (%s) for cluster %v.", machine.Name, cluster.Name)
		_, err := a.clientset.CoreV1().Pods(cluster.Namespace).Create(a.getPodSpec(cluster, machine))
		if err != nil {
			log.Printf("Creating pod (%s) for cluster %v failed: %v.", machine.Name, cluster.Name, err)
			return fmt.Errorf("Could not create the pod (%s) for cluster %s: %v", machine.Name, cluster.Name, err)
		}
	}
	// TODO: handle update pod
	return nil
}

func (a *Actuator) getPodSpec(cluster *clusterv1.Cluster, machine *clusterv1.Machine) *v1.Pod {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: machine.Name,
			Annotations: map[string]string{
				"kubernetes.io/target-runtime": "virtlet.cloud",
				"VirtletRootVolumeSize":        "8Gi",
				"VirtletVCPUCount":             "2",
				"VirtletCloudInitUserData":     a.getUserData(cluster, machine),
			},
		},
		Spec: v1.PodSpec{
			NodeSelector: map[string]string{"extraRuntime": "virtlet"},
			Containers: []v1.Container{
				{
					Name:            "k8s-node",
					Image:           "virtlet.cloud/cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img",
					ImagePullPolicy: v1.PullIfNotPresent,
					// for kubectl attach to work
					TTY:   true,
					Stdin: true,
				},
			},
		},
	}
	if util.IsMaster(machine) {
		if pod.Labels == nil {
			pod.Labels = map[string]string{}
		}
		pod.Labels["role"] = "k8s-master"
	}
	return pod
}

func (a *Actuator) getUserData(cluster *clusterv1.Cluster, machine *clusterv1.Machine) string {
	var kubeadm string
	if util.IsMaster(machine) {
		kubeadm = `
    # master node
    kubeadm init --token adcb82.4eae29627dc4c5a6 --pod-network-cidr=10.200.0.0/16 --service-cidr=10.97.0.0/16 --apiserver-cert-extra-sans=127.0.0.1,localhost
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
    sed -i "s/--insecure-port=0/--insecure-port=8080\\n    - --insecure-bind-address=0.0.0.0/" /etc/kubernetes/manifests/kube-apiserver.yaml
    pip install kubernetes
    echo "Master setup complete." >&2
` + a.getCephFiles(cluster)
	} else {
		kubeadm = `
    # worker node
    kubeadm join --token adcb82.4eae29627dc4c5a6 --discovery-token-unsafe-skip-ca-verification k8s-master:6443
    echo "Node setup complete." >&2
`
	}

	return `write_files:
- path: /etc/systemd/system/docker.service.d/env.conf
  permissions: "0644"
  owner: root
  content: |
    [Service]
    Environment="DOCKER_OPTS=--storage-driver=overlay2"
- path: /etc/apt/sources.list.d/kubernetes.list
  permissions: "0644"
  owner: root
  content: |
    deb http://apt.kubernetes.io/ kubernetes-xenial main
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
  ` + kubeadm + `
users:
- name: root
  # VirtletSSHKeys only affects 'ubuntu' user for this image, but we want root access
  ssh-authorized-keys:
  - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCaJEcFDXEK2ZbX0ZLS1EIYFZRbDAcRfuVjpstSc0De8+sV1aiu+dePxdkuDRwqFtCyk6dEZkssjOkBXtri00MECLkir6FcH3kKOJtbJ6vy3uaJc9w1ERo+wyl6SkAh/+JTJkp7QRXj8oylW5E20LsbnA/dIwWzAF51PPwF7A7FtNg9DnwPqMkxFo1Th/buOMKbP5ZA1mmNNtmzbMpMfJATvVyiv3ccsSJKOiyQr6UG+j7sc/7jMVz5Xk34Vd0l8GwcB0334MchHckmqDB142h/NCWTr8oLakDNvkfC1YneAfAO41hDkUbxPtVBG5M/o7P4fxoqiHEX+ZLfRxDtHB53 me@localhost
runcmd:
- /usr/local/bin/provision.sh
`
}

func (a *Actuator) getCephFiles(cluster *clusterv1.Cluster) string {
	provider, err := providerv1.ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return ""
	}
	adminKey := provider.CephAdminKey
	clientKey := provider.CephClientKey

	svcs, _ := a.clientset.CoreV1().Services("rook-ceph").List(metav1.ListOptions{LabelSelector: "app=rook-ceph-mon"})
	// TODO: Don't ignore error!
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

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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"

	"sigs.k8s.io/cluster-api-provider-virtlet/pkg/userdata"
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
	Kubeconfig     *rest.Config
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	var config *rest.Config
	var err error

	if params.Kubeconfig == nil {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config = params.Kubeconfig
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
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceMemory: resource.MustParse("1Gi"),
						},
					},
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
	if util.IsMaster(machine) {
		return userdata.MasterNodeUserData(a.clientset, cluster)
	}
	return userdata.ComputeNodeUserData()
}

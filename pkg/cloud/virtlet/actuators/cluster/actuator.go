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

package cluster

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
)

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	clustersGetter client.ClustersGetter
	clientset      *kubernetes.Clientset
	rookClientset  *rookclient.Clientset
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClustersGetter client.ClustersGetter
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

	rookClientset, err := rookclient.NewForConfig(config)

	return &Actuator{
		clustersGetter: params.ClustersGetter,
		clientset:      clientset,
		rookClientset:  rookClientset,
	}, nil
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	log.Printf("Reconciling cluster %v.", cluster.Name)

	a.reconcileMasterService(cluster)

	a.reconcileCephPool(cluster)

	// TODO: Craete an ingress resource(One for all clusters so only on LB IP will be used)
	// TODO: Add ingress rules for cluster service (maybe including https)

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	log.Printf("Deleting cluster %v.", cluster.Name)

	// TODO: Delete ceph pool
	// TODO: Delete service cluster
	// TODO: Remove ingress rules

	// or, just delete namespace?

	return nil
}

func (a *Actuator) reconcileMasterService(cluster *clusterv1.Cluster) error {
	_, err := a.clientset.CoreV1().Services(cluster.Namespace).Get("k8s-master", metav1.GetOptions{})
	if err != nil {
		log.Printf("Creating service 'master' for cluster %v.", cluster.Name)
		_, err := a.clientset.CoreV1().Services(cluster.Namespace).Create(getMasterServiceSpec())
		if err != nil {
			log.Printf("Creating service 'master' for cluster %v failed: %v.", cluster.Name, err)
			return fmt.Errorf("Could not create the service 'master' for cluster %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func getMasterServiceSpec() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "k8s-master",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"role": "k8s-master"},
			Ports: []v1.ServicePort{
				{Port: 6443, TargetPort: intstr.FromInt(6443), Name: "secured-api"},
				{Port: 8080, TargetPort: intstr.FromInt(8080), Name: "insecure-api"},
			},
			Type: v1.ServiceTypeLoadBalancer,
		},
	}
}

func (a *Actuator) reconcileCephPool(cluster *clusterv1.Cluster) error {
	_, err := a.rookClientset.CephV1().CephBlockPools("rook-ceph").Get(cluster.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Creating Ceph pool for cluster %v.", cluster.Name)
		_, err := a.rookClientset.CephV1().CephBlockPools("rook-ceph").Create(getCephPoolSpec(cluster.Name))
		if err != nil {
			log.Printf("Creating Ceph Pool for cluster %v failed: %v.", cluster.Name, err)
			return fmt.Errorf("Could not create Ceph pool for cluster %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func getCephPoolSpec(name string) *cephv1.CephBlockPool {
	return &cephv1.CephBlockPool{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: cephv1.PoolSpec{
			FailureDomain: "host",
			Replicated:    cephv1.ReplicatedSpec{Size: 1},
		},
	}
}

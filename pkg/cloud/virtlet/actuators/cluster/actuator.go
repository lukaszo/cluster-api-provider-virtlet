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
	"strings"

	v1 "k8s.io/api/core/v1"
	extensionsbeta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	providerv1 "sigs.k8s.io/cluster-api-provider-virtlet/pkg/apis/virtlet/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
)

const (
	PROVIDER_RBAC_NAME = "cluster-api-provider-virtlet-rbac"
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

	err := a.reconcileAPIServerService(cluster)
	if err != nil {
		return fmt.Errorf("Error when reconciling API Server service: %v", err)
	}

	err = a.reconcileCephPool(cluster)
	if err != nil {
		return fmt.Errorf("Error when reconciling ceph pool: %v", err)
	}

	err = a.reconcileIngress(cluster)
	if err != nil {
		return fmt.Errorf("Error when reconciling ingress: %v", err)
	}

	err = a.reconcileRBAC(cluster)
	if err != nil {
		return fmt.Errorf("Error when reconciling RBAC: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	log.Printf("Deleting cluster %v.", cluster.Name)
	var err error

	cErr := a.deleteCephPool(cluster)
	if err != nil {
		log.Printf("Error when deleting ceph pool for cluster %s: %v", cluster.Name, err)
	}

	sErr := a.deleteAPIServerService(cluster)
	if err != nil {
		log.Printf("Error when deleting APIServer service for cluster %s: %v", cluster.Name, err)
	}

	iErr := a.deleteIngress(cluster)
	if err != nil {
		log.Printf("Error when deleting ingress for cluster %s: %v", cluster.Name, err)
	}

	if cErr != nil || sErr != nil || iErr != nil {
		return fmt.Errorf("Cluster %s delete failed", cluster.Name)
	}

	// TODO: delete RBAC rules

	return nil
}

func (a *Actuator) reconcileAPIServerService(cluster *clusterv1.Cluster) error {
	_, err := a.clientset.CoreV1().Services(cluster.Namespace).Get("api-server", metav1.GetOptions{})
	if err != nil {
		_, err := a.clientset.CoreV1().Services(cluster.Namespace).Create(getAPIServerServiceSpec())
		if err != nil {
			return fmt.Errorf("Could not create the service 'api-server' for cluster %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func (a *Actuator) deleteAPIServerService(cluster *clusterv1.Cluster) error {
	_, err := a.clientset.CoreV1().Services(cluster.Namespace).Get("api-server", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
	}
	err = a.clientset.CoreV1().Services(cluster.Namespace).Delete("api-server", &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Could not delete the service 'api-server' for cluster %s: %v", cluster.Name, err)
	}
	return nil
}

func getAPIServerServiceSpec() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "api-server",
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
		_, err := a.rookClientset.CephV1().CephBlockPools("rook-ceph").Create(getCephPoolSpec(cluster.Name))
		if err != nil {
			return fmt.Errorf("Could not create Ceph pool for cluster %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func (a *Actuator) deleteCephPool(cluster *clusterv1.Cluster) error {
	_, err := a.rookClientset.CephV1().CephBlockPools("rook-ceph").Get(cluster.Name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
	}
	err = a.rookClientset.CephV1().CephBlockPools("rook-ceph").Delete(cluster.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Could not delete the ceph pool %s for cluster %s: %v", cluster.Name, cluster.Name, err)
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

func (a *Actuator) reconcileIngress(cluster *clusterv1.Cluster) error {
	providerConf, err := providerv1.ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return fmt.Errorf("Couldn't get cluster providerSpec for cluster %s: %v", cluster.Name, err)
	}
	// TODO: check if this is a valid host
	host := providerConf.Host

	_, err = a.clientset.ExtensionsV1beta1().Ingresses(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
	if err != nil {
		_, err := a.clientset.ExtensionsV1beta1().Ingresses(cluster.Namespace).Create(getIngressSpec(cluster.Name, host))
		if err != nil {
			return fmt.Errorf("Could not create an Ingress for cluster %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func (a *Actuator) deleteIngress(cluster *clusterv1.Cluster) error {
	_, err := a.clientset.ExtensionsV1beta1().Ingresses(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
	}
	err = a.clientset.ExtensionsV1beta1().Ingresses(cluster.Namespace).Delete(cluster.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Could not delete the ingress %s for cluster %s: %v", cluster.Name, cluster.Name, err)
	}
	return nil
}

func getIngressSpec(name, host string) *extensionsbeta1.Ingress {
	ingress := &extensionsbeta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: extensionsbeta1.IngressSpec{
			Rules: []extensionsbeta1.IngressRule{
				{
					Host: host,
				},
			},
		},
	}
	ingress.Spec.Rules[0].HTTP = &extensionsbeta1.HTTPIngressRuleValue{
		Paths: []extensionsbeta1.HTTPIngressPath{
			{
				Path: "/",
				Backend: extensionsbeta1.IngressBackend{
					ServiceName: "api-server",
					ServicePort: intstr.FromInt(8080),
				},
			},
		},
	}
	return ingress
}

func (a *Actuator) reconcileRBAC(cluster *clusterv1.Cluster) error {
	err := a.reconcileRoles(cluster)
	if err != nil {
		return fmt.Errorf("Couldn't reconcile RBAC rules for cluster %s: %v", cluster.Name, err)
	}
	err = a.reconcileRoleBindings(cluster)
	if err != nil {
		return fmt.Errorf("Couldn't reconcile RBAC rules for cluster %s: %v", cluster.Name, err)
	}
	return nil
}

func (a *Actuator) reconcileRoles(cluster *clusterv1.Cluster) error {
	_, err := a.clientset.RbacV1().ClusterRoles().Get(PROVIDER_RBAC_NAME, metav1.GetOptions{})
	if err != nil {
		_, err := a.clientset.RbacV1().ClusterRoles().Create(getRBACRoleSpec())
		if err != nil {
			return fmt.Errorf("Could not create the Role %s for cluster %s: %v", PROVIDER_RBAC_NAME, cluster.Name, err)
		}
	}
	return nil
}

func (a *Actuator) reconcileRoleBindings(cluster *clusterv1.Cluster) error {
	// Binding
	_, err := a.clientset.RbacV1().ClusterRoleBindings().Get(PROVIDER_RBAC_NAME, metav1.GetOptions{})
	if err != nil {
		roleBinding := getRoleBindingsSpec("default", cluster.Namespace, "ServiceAccount",
			"rbac.authorization.k8s.io", "ClusterRole", PROVIDER_RBAC_NAME)
		_, err := a.clientset.RbacV1().ClusterRoleBindings().Create(roleBinding)
		if err != nil {
			return fmt.Errorf("Could not create the RoleBinding %s for cluster %s, binding %v: %v", PROVIDER_RBAC_NAME, cluster.Name, roleBinding, err)
		}
	}
	return nil
}

func getRBACRoleSpec() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: PROVIDER_RBAC_NAME,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"virtletlb.virtlet.cloud"},
				Resources: []string{"innerservices"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"*"},
			},
		},
	}
}

func getRoleBindingsSpec(subjectName, subjectNamespace, subjectKind, roleAPIGroup, roleKind, roleName string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PROVIDER_RBAC_NAME,
			Namespace: subjectNamespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      subjectKind,
				Name:      subjectName,
				Namespace: subjectNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: roleAPIGroup,
			Kind:     roleKind,
			Name:     roleName,
		},
	}
}

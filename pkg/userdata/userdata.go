package userdata

import (
	"k8s.io/client-go/kubernetes"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// TODO: check out https://github.com/kubermatic/machine-controller/tree/master/pkg/userdata
func MasterNodeUserData(clientset *kubernetes.Clientset, cluster *clusterv1.Cluster) string {
	return `
write_files:` + masterProvisionScript + getCephFiles(clientset, cluster) + getKubernetesAPTKey() + getDockerService() + `
users: ` + getUsers() + `
runcmd:
- /usr/local/bin/provision.sh
`
}

func ComputeNodeUserData() string {
	return `
write_files:` + computeProvisionScript + getKubernetesAPTKey() + getDockerService() + `
users: ` + getUsers() + `
runcmd:
- /usr/local/bin/provision.sh
`
}

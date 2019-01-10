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
package controller

import (
	"sigs.k8s.io/cluster-api-provider-virtlet/pkg/cloud/virtlet/actuators/machine"
	capimachine "sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

//+kubebuilder:rbac:groups=virtlet.k8s.io,resources=virtletmachineproviderspecs;virtletmachineproviderstatuses,verbs=get;list;watch;create;update;patch;delete
func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, func(m manager.Manager) error {
		return capimachine.AddWithActuator(m, &machine.Actuator{})
	})
}

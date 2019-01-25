// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletClusterProviderSpec) DeepCopyInto(out *VirtletClusterProviderSpec) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletClusterProviderSpec.
func (in *VirtletClusterProviderSpec) DeepCopy() *VirtletClusterProviderSpec {
	if in == nil {
		return nil
	}
	out := new(VirtletClusterProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletClusterProviderSpec) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletClusterProviderStatus) DeepCopyInto(out *VirtletClusterProviderStatus) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletClusterProviderStatus.
func (in *VirtletClusterProviderStatus) DeepCopy() *VirtletClusterProviderStatus {
	if in == nil {
		return nil
	}
	out := new(VirtletClusterProviderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletClusterProviderStatus) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletClusterProviderStatusList) DeepCopyInto(out *VirtletClusterProviderStatusList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtletClusterProviderStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletClusterProviderStatusList.
func (in *VirtletClusterProviderStatusList) DeepCopy() *VirtletClusterProviderStatusList {
	if in == nil {
		return nil
	}
	out := new(VirtletClusterProviderStatusList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletClusterProviderStatusList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletClusterProviderStatusSpec) DeepCopyInto(out *VirtletClusterProviderStatusSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletClusterProviderStatusSpec.
func (in *VirtletClusterProviderStatusSpec) DeepCopy() *VirtletClusterProviderStatusSpec {
	if in == nil {
		return nil
	}
	out := new(VirtletClusterProviderStatusSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletClusterProviderStatusStatus) DeepCopyInto(out *VirtletClusterProviderStatusStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletClusterProviderStatusStatus.
func (in *VirtletClusterProviderStatusStatus) DeepCopy() *VirtletClusterProviderStatusStatus {
	if in == nil {
		return nil
	}
	out := new(VirtletClusterProviderStatusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderSpec) DeepCopyInto(out *VirtletMachineProviderSpec) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderSpec.
func (in *VirtletMachineProviderSpec) DeepCopy() *VirtletMachineProviderSpec {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletMachineProviderSpec) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderSpecList) DeepCopyInto(out *VirtletMachineProviderSpecList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtletMachineProviderSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderSpecList.
func (in *VirtletMachineProviderSpecList) DeepCopy() *VirtletMachineProviderSpecList {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderSpecList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletMachineProviderSpecList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderSpecSpec) DeepCopyInto(out *VirtletMachineProviderSpecSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderSpecSpec.
func (in *VirtletMachineProviderSpecSpec) DeepCopy() *VirtletMachineProviderSpecSpec {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderSpecSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderSpecStatus) DeepCopyInto(out *VirtletMachineProviderSpecStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderSpecStatus.
func (in *VirtletMachineProviderSpecStatus) DeepCopy() *VirtletMachineProviderSpecStatus {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderSpecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderStatus) DeepCopyInto(out *VirtletMachineProviderStatus) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderStatus.
func (in *VirtletMachineProviderStatus) DeepCopy() *VirtletMachineProviderStatus {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletMachineProviderStatus) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderStatusList) DeepCopyInto(out *VirtletMachineProviderStatusList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtletMachineProviderStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderStatusList.
func (in *VirtletMachineProviderStatusList) DeepCopy() *VirtletMachineProviderStatusList {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderStatusList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtletMachineProviderStatusList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderStatusSpec) DeepCopyInto(out *VirtletMachineProviderStatusSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderStatusSpec.
func (in *VirtletMachineProviderStatusSpec) DeepCopy() *VirtletMachineProviderStatusSpec {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderStatusSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtletMachineProviderStatusStatus) DeepCopyInto(out *VirtletMachineProviderStatusStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtletMachineProviderStatusStatus.
func (in *VirtletMachineProviderStatusStatus) DeepCopy() *VirtletMachineProviderStatusStatus {
	if in == nil {
		return nil
	}
	out := new(VirtletMachineProviderStatusStatus)
	in.DeepCopyInto(out)
	return out
}

// Copyright 2019 IBM Corp
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ImageDirDefault = "/var/lib/kubelet/secure-images"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecureContainerImageConfigSpec defines the desired state of SecureContainerImageConfig
// +k8s:openapi-gen=true
type SecureContainerImageConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ImageDir         string `json:"imageDir"`
	RuntimeClassName string `json:"runtimeClassName"`
}

// SecureContainerImageConfigStatus defines the observed state of SecureContainerImageConfig
// +k8s:openapi-gen=true
type SecureContainerImageConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainerImageConfig is the Schema for the securecontainerimageconfigs API
// +k8s:openapi-gen=true
type SecureContainerImageConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecureContainerImageConfigSpec   `json:"spec,omitempty"`
	Status SecureContainerImageConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainerImageConfigList contains a list of SecureContainerImageConfig
type SecureContainerImageConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecureContainerImageConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecureContainerImageConfig{}, &SecureContainerImageConfigList{})
}

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
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecureContainerSpec defines the desired state of SecureContainer
// +k8s:openapi-gen=true
type SecureContainerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	SecureContainerImageRef SecureContainerImageRef `json:"SecureContainerImageRef"`
}

// SecureContainerImageRef defines the SecureContainerImage of SecureContainer
// +k8s:openapi-gen=true
type SecureContainerImageRef struct {
	Name string `json:"name"`
}

// SecureContainerStatus defines the observed state of SecureContainer
// +k8s:openapi-gen=true
type SecureContainerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainer is the Schema for the securecontainers API
// +k8s:openapi-gen=true
type SecureContainer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecureContainerSpec   `json:"spec,omitempty"`
	Status SecureContainerStatus `json:"status,omitempty"`
	Object runtime.RawExtension  `json:"object"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainerList contains a list of SecureContainer
type SecureContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecureContainer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecureContainer{}, &SecureContainerList{})
}

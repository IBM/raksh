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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecureContainerImageSpec defines the desired state of SecureContainerImage
// +k8s:openapi-gen=true
type SecureContainerImageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	VMImage                        string                         `json:"vmImage"`
	ImagePullSecrets               []corev1.LocalObjectReference  `json:"imagePullSecrets,omitempty"`
	SecureContainerImageConfigRef  SecureContainerImageConfigRef  `json:"SecureContainerImageConfigRef"`
	SecureContainerImageConfigSpec SecureContainerImageConfigSpec `json:"SecureContainerImageConfigSpec,omitempty"`
}

// SecureContainerImageConfigRef defines SecureContainerImage configuration
// +k8s:openapi-gen=true
type SecureContainerImageConfigRef struct {
	Name string `json:"name"`
}

// SecureContainerImageStatus defines the observed state of SecureContainerImage
// +k8s:openapi-gen=true
type SecureContainerImageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainerImage is the Schema for the securecontainerimages API
// +k8s:openapi-gen=true
type SecureContainerImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecureContainerImageSpec   `json:"spec,omitempty"`
	Status SecureContainerImageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecureContainerImageList contains a list of SecureContainerImage
type SecureContainerImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecureContainerImage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecureContainerImage{}, &SecureContainerImageList{})
}

func (i *SecureContainerImage) GetImageDirWithName() (image string, err error) {
	if i.Spec.SecureContainerImageConfigSpec.ImageDir == "" {
		return "", fmt.Errorf("Emtpy ImageDir")
	}
	return i.Spec.SecureContainerImageConfigSpec.ImageDir + "/" + i.Name, nil
}

func (i *SecureContainerImage) GetRuntimeClassName() string {
	return i.Spec.SecureContainerImageConfigSpec.RuntimeClassName
}

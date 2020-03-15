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

package app

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestMountRakshSecrets(t *testing.T) {
	testRakshSecretName := "test-raksh-secret"
	testVolumeName := "secure-volume-raksh"

	expectedVolume := corev1.Volume{
		Name: testVolumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: testRakshSecretName,
			},
		},
	}

	expectedVolmount := corev1.VolumeMount{
		Name:      testVolumeName,
		ReadOnly:  true,
		MountPath: "/etc/raksh-secrets",
	}

	testpod := &corev1.PodSpec{
		Containers: []corev1.Container{
			{Name: "c1"},
		},
	}

	mountRakshSecrets(testpod, testRakshSecretName)

	for i := range testpod.Containers {
		if equal := reflect.DeepEqual(expectedVolmount, testpod.Containers[i].VolumeMounts[0]); !equal {
			t.Fatalf("For container %s, actual : %+v is not matching the expected: %+v", testpod.Containers[i].Name, testpod.Containers[i].VolumeMounts[0], expectedVolmount)
		}
		if equal := reflect.DeepEqual(expectedVolume, testpod.Volumes[0]); !equal {
			t.Fatalf("For container %s, actual : %+v is not matching the expected: %+v", testpod.Containers[i].Name, testpod.Volumes[0], expectedVolume)
		}
	}
}

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

package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ibm/raksh/pkg/crypto"
	"github.com/ibm/raksh/tests/rakshctl/framework"
)

var sampleWorkload = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
`

func TestImageAppCreate(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, randImageName(),
		"--filename", sampleWorkloadYaml, "--image", "sample-securecontainerimage", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	// Validation block - open the securecontainer file and check the content
	content, err := ioutil.ReadFile(filepath.Join(tc.dir, "sample-sc.yaml"))
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}
	if !strings.Contains(string(content), "kind: SecureContainer") {
		t.Errorf("\nExpecting the \"kind: SecureContainer\" in the generated file but the actual content is: \n%s", string(content))
	}
}

// Ensure that buildCmd been used to create the secureimage
func TestBuildCmd(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	dir, err := ioutil.TempDir("", "TestBuildCmd")
	if err != nil {
		t.Fatalf("Error to create tempdir: %+v", err)
	}
	defer os.RemoveAll(dir)

	customBuildCmd := filepath.Join(dir, "custom-docker")
	defaultBuildCMDPath, err := exec.LookPath(framework.DefaultBuildCmd)
	if err != nil {
		t.Fatalf("Failed to locate the DefaultBuildCmd: %s", framework.DefaultBuildCmd)
	}
	os.Symlink(defaultBuildCMDPath, customBuildCmd)

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	expectedSecureImage := randImageName()

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, expectedSecureImage,
		"--filename", sampleWorkloadYaml, "--image", "sample-securecontainerimage", "--buildCmd", customBuildCmd, "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}
	// TODO: add more validation later
}

// Ensure command to fail if wrong buildCMD mentioned
func TestInvalidBuildCmd(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	dir, err := ioutil.TempDir("", "TestBuildCmd")
	if err != nil {
		t.Fatalf("Error to create tempdir: %+v", err)
	}
	defer os.RemoveAll(dir)

	customBuildCmd := filepath.Join(dir, "custom-docker")

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	expectedSecureImage := randImageName()

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, expectedSecureImage, "--filename", sampleWorkloadYaml,
		"--image", "sample-securecontainerimage", "--buildCmd", customBuildCmd}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	if err == nil {
		t.Errorf("Test passed instead of failing, stdout: %+v, stderr: %+v", std, stderr)
	}
	// TODO: add more validation later
}

func TestInsecureOperation(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, randImageName(),
		"--filename", sampleWorkloadYaml, "--image", "sample-securecontainerimage", "--symmKeyFile", tc.keyFile,
		"--nonceFile", tc.nonceFile, "--output", tc.dir, "--insecure"}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	// Validation block - open the raksh-secret.yaml and check the content
	// The secrets should be same as provided via input
	content, err := ioutil.ReadFile(filepath.Join(tc.dir, "raksh-secret.yaml"))
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	keyEnc, nonceEnc, err := crypto.GetConfigMapKeys(tc.keyFile, tc.nonceFile)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	//Check configMapKey
	if !strings.Contains(string(content), string(keyEnc)) {
		t.Errorf("\nExpecting key: %s in the generated file but the actual content is: \n%s", string(keyEnc), string(content))
	}
	//Check nonce
	if !strings.Contains(string(content), string(nonceEnc)) {
		t.Errorf("\nExpecting nonce: %s in the generated file but the actual content is: \n%s", string(nonceEnc), string(content))
	}

}

func TestSecureOperation(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, randImageName(),
		"--filename", sampleWorkloadYaml, "--image", "sample-securecontainerimage", "--symmKeyFile", tc.keyFile,
		"--nonceFile", tc.nonceFile, "--output", tc.dir}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	// Validation block - open the raksh-secret.yaml and check the content
	// The secrets should be same as provided via input
	content, err := ioutil.ReadFile(filepath.Join(tc.dir, "raksh-secret.yaml"))
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	//Check configMapKey. It should be rakshDummyConfigMapKey
	if !strings.Contains(string(content), "rakshDummyConfigMapKey") {
		t.Errorf("\nExpecting key: rakshDummyConfigMapKey in the generated file but the actual content is: \n%s", string(content))
	}
	//Check nonce. It should rakshDummyNonce
	if !strings.Contains(string(content), "rakshDummyNonce") {
		t.Errorf("\nExpecting nonce: rakshDummyNonce in the generated file but the actual content is: \n%s", string(content))
	}

}

func TestRakshSecretFileGeneration(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	sampleWorkloadYaml := filepath.Join(tc.dir, "sample.yaml")
	fd, err := os.Create(sampleWorkloadYaml)
	if err != nil {
		t.Fatalf("Error to create %s file", sampleWorkloadYaml)
	}
	defer fd.Close()
	fd.WriteString(sampleWorkload)

	var cmdArgs = []string{"image", "create", "--initrd", tc.initrd, "--vmlinux", tc.vmlinux, randImageName(),
		"--filename", sampleWorkloadYaml, "--image", "sample-securecontainerimage", "--symmKeyFile", tc.keyFile,
		"--nonceFile", tc.nonceFile}

	std, stderr, err := executeCommand("rakshctl", cmdArgs...)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

	//Validate block - Check if raksh-secret.yaml got generated or not
	_, err = os.Stat("raksh-secret.yaml")
	if err != nil {
		t.Errorf("\nGot the error: %+v", err)
	}

}

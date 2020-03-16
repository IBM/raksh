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
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	randutil "github.com/ibm/raksh/pkg/utils/random"

	b64 "encoding/base64"

	"github.com/ibm/raksh/pkg/rakshctl/cmd/image"
	"github.com/ibm/raksh/pkg/utils/cmd"
)

type testcontext struct {
	dir       string
	initrd    string
	vmlinux   string
	imageName string
	keyFile   string
	nonceFile string
}

var imageCreateCMD = []string{"image", "create"}

func executeCommand(command string, args ...string) (output string, stderr string, err error) {
	output, stderr, err = cmd.Exec(command, args)
	return output, stderr, err
}

// randImageName returns the random imagename in the format of securecontainerimage-<5 random images>:latest
func randImageName() string {
	chars := make([]byte, 5)
	for i := 0; i < 5; i++ {
		chars[i] = byte(97 + rand.Intn(25))
	}
	return "securecontainerimage-" + string(chars) + ":latest"
}

func (tc *testcontext) setup(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestExactValidArgs")
	if err != nil {
		t.Fatalf("Error to create tempdir: %+v", err)
	}
	tc.dir = dir

	initrdFile := filepath.Join(dir, "initrd")
	vmlinuxFile := filepath.Join(dir, "vmlinux")

	initrdFD, err := os.Create(initrdFile)
	if err != nil {
		t.Fatalf("Error to create %s file", initrdFile)
	}
	defer initrdFD.Close()
	tc.initrd = initrdFile

	vmlinuxFD, err := os.Create(vmlinuxFile)
	if err != nil {
		t.Fatalf("Error to create %s file", initrdFile)
	}
	defer vmlinuxFD.Close()

	symmKeyFile := filepath.Join(dir, "symmKey")

	keyFD, err := os.Create(symmKeyFile)
	if err != nil {
		t.Fatalf("Error to create symmKey file: %+v", err)
	}

	buf, err := randutil.GetBytes(32)
	if err != nil {
		t.Fatalf("Unable to get random bytes for key: %+v", err)
		keyFD.Close()
	}

	//Encode to base64 and write the key to the file
	bufEncStr := b64.StdEncoding.EncodeToString(buf)
	if err != nil {
		t.Fatalf("Unable to encode random bytes for key: %+v", err)
		keyFD.Close()
	}

	_, err = keyFD.Write([]byte(bufEncStr))
	if err != nil {
		t.Fatalf("Error in wirting key to symmKey file: %+v", err)
		keyFD.Close()
	}

	nonceFile := filepath.Join(dir, "nonce")

	nonceFD, err := os.Create(nonceFile)
	if err != nil {
		t.Fatalf("Error to create nonce file: %+v", err)
	}

	buf, err = randutil.GetBytes(12)
	if err != nil {
		t.Fatalf("Unable to get random bytes for nonce: %+v", err)
		nonceFD.Close()
	}

	//Encode to base64 and write the key to the file
	bufEncStr = b64.StdEncoding.EncodeToString(buf)
	if err != nil {
		t.Fatalf("Unable to encode random bytes for key: %+v", err)
		nonceFD.Close()
	}

	_, err = nonceFD.Write([]byte(bufEncStr))
	if err != nil {
		t.Fatalf("Error in wirting key to nonce file: %+v", err)
		nonceFD.Close()
	}

	tc.vmlinux = vmlinuxFile
	tc.imageName = randImageName()
	tc.keyFile = symmKeyFile
	tc.nonceFile = nonceFile
}

func (tc *testcontext) teardown() {
	os.RemoveAll(tc.dir)
}

func TestWithDifferentArgs(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	var createTests = []struct {
		name string
		args []string
		err  string
	}{
		{
			"Positive test",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			"",
		},
		{
			"Missing vmlinux argument",
			[]string{"--initrd", tc.initrd, "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			`required flag(s) "vmlinux" not set`,
		},
		{
			"Missing vmlinux and initrd arguments",
			[]string{"--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			`required flag(s) "initrd", "vmlinux" not set`,
		},
		{
			"Non-existent initrd",
			[]string{"--initrd", "fakeinitrd", "--vmlinux", tc.vmlinux, "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			`open fakeinitrd: no such file or directory`,
		},
		{
			"Non-existent vmlinux",
			[]string{"--initrd", tc.initrd, "--vmlinux", "fakevmlinux", "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			`Failed to copy fakevmlinux to`,
		},
		{
			"Fail docker push",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "--push", "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			"Failed to push the",
		},
		{
			"Non-existing symmKey",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "--skip-app", "--nonceFile", tc.nonceFile},
			`required flag(s) "symmKeyFile" not set`,
		},
		{
			"Non-existing nonce",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "--skip-app", "--symmKeyFile", tc.keyFile},
			`required flag(s) "nonceFile" not set`,
		},
	}
	for _, tt := range createTests {
		fmt.Println("Running test: ", tt.name)
		fmt.Println("Args: ", tt.args)
		imageName := randImageName()
		defer cmd.Exec("docker", []string{"rmi", "-f", imageName})

		tt.args = append(imageCreateCMD, tt.args...)
		tt.args = append(tt.args, imageName)

		std, stderr, err := executeCommand("rakshctl", tt.args...)
		fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
		if tt.err != "" && !strings.Contains(stderr, tt.err) {
			t.Errorf("\ntest: %s, \nexpected: %+v \nbut, got: %+v", tt.name, tt.err, stderr)
		}
	}
}

func TestImageName(t *testing.T) {
	tc := &testcontext{}
	tc.setup(t)
	defer tc.teardown()

	var createTests = []struct {
		name string
		args []string
		err  string
	}{
		{
			"Missing Image Name",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			image.MissingImageError,
		},
		{
			"Invalid Image Name",
			[]string{"--initrd", tc.initrd, "--vmlinux", tc.vmlinux, "invalid-IMAGENAME:latest", "--skip-app", "--symmKeyFile", tc.keyFile, "--nonceFile", tc.nonceFile},
			"Failed to build the SecureContainer image",
		},
	}
	for _, tt := range createTests {
		tt.args = append(imageCreateCMD, tt.args...)

		_, stderr, err := executeCommand("rakshctl", tt.args...)
		if tt.err != "" && !strings.Contains(stderr, tt.err) {
			t.Errorf("\ntest: %s, \nexpected: %+v \nbut, got: %+v", tt.name, tt.err, err)
		}
	}
}

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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	randutil "github.com/ibm/raksh/pkg/utils/random"

	"github.com/ibm/raksh/pkg/rakshctl/cmd/app"
	"github.com/ibm/raksh/pkg/utils/cmd"
)

var unsupportedWorkload = `
apiVersion: v1
kind: Service
metadata:
  name: redis-master
  labels:
    app: redis
    role: master
    tier: backend
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    app: redis
    role: master
    tier: backend
`

func TestAppCreateUnsupportedWorkload(t *testing.T) {
	dir, err := ioutil.TempDir("", "AppCreate")
	if err != nil {
		t.Fatalf("Error to create tempdir: %+v", err)
	}
	defer os.RemoveAll(dir)

	workload := filepath.Join(dir, "workload.yaml")
	workloadFD, err := os.Create(workload)
	if err != nil {
		t.Fatalf("Error to create %s file", workload)
	}
	defer workloadFD.Close()
	_, err = workloadFD.WriteString(unsupportedWorkload)
	if err != nil {
		t.Fatalf("Failed to write to %s file", workload)
	}

	symmKeyFile := dir + "/symm_key"

	f, err := os.Create(symmKeyFile)
	if err != nil {
		t.Fatalf("Error to create symmKeyFile file: %+v", err)
	}

	buf, err := randutil.GetBytes(32)
	if err != nil {
		t.Fatalf("Unable to get random bytes for key: %+v", err)
		f.Close()
	}

	_, err = f.Write(buf)
	if err != nil {
		t.Fatalf("Error in wirting key to symmKeyFile file: %+v", err)
		f.Close()
	}

	var cmdArgs = []string{"app", "create", "-f", workload, "-i", "securecontainerimage-example", "--symmKeyFile", symmKeyFile}
	std, stderr, err := cmd.Exec("rakshctl", cmdArgs)
	fmt.Printf("stdout: %+v, stderr: %+v, err: %+v", std, stderr, err)
	exp := fmt.Sprintf(app.UnsupportedKindMsg, workload)
	if !strings.Contains(std, exp) {
		t.Errorf("\nExpected to contain %s in \"%s\"", exp, std)
	}
}

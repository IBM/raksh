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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ibm/raksh/pkg/crypto"
	"github.com/ibm/raksh/pkg/rakshctl/types/flags"
	cpioutil "github.com/ibm/raksh/pkg/utils/cpio"
	gziputil "github.com/ibm/raksh/pkg/utils/gzip"
	randutil "github.com/ibm/raksh/pkg/utils/random"
)

func TestAppendFilesToInitrd(t *testing.T) {
	dir, err := ioutil.TempDir("", "image-tests")
	if err != nil {
		t.Fatalf("Error to create tempdir: %+v", err)
	}

	symmKeyFile := dir + "/symm_key"
	flags.Key = symmKeyFile

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

	defer os.RemoveAll(dir)

	var fileset1 = []cpioutil.File{
		{Name: "fileset1_file1", Body: randomBytes()},
		{Name: "fileset1_file2", Body: randomBytes()},
	}
	initrd1 := filepath.Join(dir, "initrd1")

	if err := createGzipFile(initrd1, fileset1); err != nil {
		t.Fatalf("Error to test fileset1 Gzipfile: %+v", err)
	}

	resultInitrd := filepath.Join(dir, "result")

	err = appendKeys(initrd1, resultInitrd)
	if err != nil {
		t.Fatalf("Failed to append the keys to initrd: %+v", err)
	}

	symmKey, nonce, _ := crypto.GetConfigMapKeys(symmKeyFile)
	var keys = []cpioutil.File{
		{Name: "symm_key", Body: symmKey},
		{Name: "nonce", Body: nonce},
	}

	actualFiles, err := extractInitrd(resultInitrd)
	if err != nil {
		t.Fatalf("Failed to extract the files from the result initrd: %+v", err)
	}

	expectedFiles := append(fileset1, keys...)
	if equal := reflect.DeepEqual(expectedFiles, actualFiles); !equal {
		t.Fatalf("actual : %+v is not matching the expected: %+v", actualFiles, expectedFiles)
	}
}

// extractGzip will open the initrd file and returns the files present
func extractInitrd(filename string) ([]cpioutil.File, error) {
	files := []cpioutil.File{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return files, err
	}

	ed := gziputil.Extract(data)

	for _, e := range ed {
		con, err := cpioutil.Extract(e)
		if err != nil {
			return files, err
		}
		for _, c := range con {
			files = append(files, c)
		}
	}
	return files, nil
}

func createGzipFile(filename string, files []cpioutil.File) error {
	cbytes, err := cpioutil.Create(files)
	if err != nil {
		return err
	}

	gbytes, err := gziputil.Create(cbytes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, gbytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func randomBytes() []byte {
	bytes, _ := randutil.GetBytes(32)
	return bytes
}

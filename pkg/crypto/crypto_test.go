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

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncryptConfigMapFiles(t *testing.T) {
	testString := "Sample String to Test Symmetric Encryption"

	dir1, err := ioutil.TempDir("", "TestEncryption-1")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir1)

	keyPath := dir1 + "symmKeyFile"
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile(keyPath, key, 0644)
	if err != nil {
		t.Error(err)
	}

	encConfigMap, err := EncryptConfigMap([]byte(testString), keyPath)

	if err != nil {
		t.Errorf("Failed to encrypt configMap %v", err)
	}

	decryptedConfigMap := decryptConfigMap(encConfigMap, keyPath)

	if testString != decryptedConfigMap {
		t.Errorf("Failed to match decryped string, expected %s but got %s", testString, decryptedConfigMap)
	}
}

func decryptConfigMap(configMap string, keyPath string) string {

	sDec, _ := b64.StdEncoding.DecodeString(configMap)

	key, nonce, err := GetConfigMapKeys(keyPath)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintextBytes, err := aesgcm.Open(nil, nonce, sDec, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintextBytes)
}

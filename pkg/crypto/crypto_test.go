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

	//Encode the Key
	keyEnc := b64.StdEncoding.EncodeToString(key)

	err = ioutil.WriteFile(keyPath, []byte(keyEnc), 0644)
	if err != nil {
		t.Error(err)
	}

	noncePath := dir1 + "nonceFile"
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		t.Error(err)
	}

	//Encode the nonce
	nonceEnc := b64.StdEncoding.EncodeToString(nonce)

	err = ioutil.WriteFile(noncePath, []byte(nonceEnc), 0644)
	if err != nil {
		t.Error(err)
	}

	encConfigMap, err := EncryptConfigMap([]byte(testString), keyPath, noncePath)

	if err != nil {
		t.Errorf("Failed to encrypt configMap %v", err)
	}

	decryptedConfigMap := decryptConfigMap(encConfigMap, keyPath, noncePath)

	if testString != decryptedConfigMap {
		t.Errorf("Failed to match decryped string, expected %s but got %s", testString, decryptedConfigMap)
	}
}

func decryptConfigMap(configMap string, keyPath string, noncePath string) string {

	sDec, _ := b64.StdEncoding.DecodeString(configMap)

	key, nonce, err := GetConfigMapKeys(keyPath, noncePath)

	//The keys are base64 encoded. Decode it before passing to encryption function
	keyDecoded, err := b64.StdEncoding.DecodeString(string(key))
	if err != nil {
		panic(err.Error())
	}
	nonceDecoded, err := b64.StdEncoding.DecodeString(string(nonce))
	if err != nil {
		panic(err.Error())
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintextBytes, err := aesgcm.Open(nil, nonceDecoded, sDec, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintextBytes)
}

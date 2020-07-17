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

	decryptedConfigMap, err := decryptConfigMap(encConfigMap, keyPath, noncePath)
	if err != nil {
		t.Errorf("Failed to decrypt configMap %v", err)
	}

	if testString != decryptedConfigMap {
		t.Errorf("Failed to match decryped string, expected %s but got %s", testString, decryptedConfigMap)
	}
}

func decryptConfigMap(configMap string, keyPath string, noncePath string) (string, error) {

	sDec, _ := b64.StdEncoding.DecodeString(configMap)

	key, nonce, err := GetConfigMapKeys(keyPath, noncePath)

	//The keys are base64 encoded. Decode it before passing to encryption function
	keyDecoded, err := b64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
	}
	nonceDecoded, err := b64.StdEncoding.DecodeString(string(nonce))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintextBytes, err := aesgcm.Open(nil, nonceDecoded, sDec, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}

func TestDecryptConfigMapWrongkey(t *testing.T) {
	/*Test Objective : A new unit test case, TestDecryptConfigMapWrongkey is being
	added for crypto. In this test, behavior of passing wrong key value
	is being tested. In present behavior, decrypt functions retuns the global var
	when value passed is non NULL.*/

	testString := "Sample String to Test Symmetric Encryption"

	dir1, err := ioutil.TempDir("", "TestEncryption-2")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir1)

	keyPath := dir1 + "symmKeyFile"
	wrongkeyPath := dir1 + "wrongsymmKeyFile"

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}
	wrongkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, wrongkey); err != nil {
		t.Error(err)
	}

	//Encode the Key
	keyEnc := b64.StdEncoding.EncodeToString(key)
	//encode wrong key
	wrongkeyEnc := b64.StdEncoding.EncodeToString(wrongkey)

	err = ioutil.WriteFile(keyPath, []byte(keyEnc), 0644)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(wrongkeyPath, []byte(wrongkeyEnc), 0644)
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

	wrongdecryptedConfigMap, err := decryptConfigMap(encConfigMap, wrongkeyPath, noncePath)
	if err != nil {
		t.Errorf("Failed to decrypt configMap %v", err)
	}

	if testString != wrongdecryptedConfigMap {
		t.Logf("Failed to match decryped string, expected %s but got %s", testString, wrongdecryptedConfigMap)
		t.FailNow()
	}
}

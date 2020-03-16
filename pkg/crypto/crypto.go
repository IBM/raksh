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
	b64 "encoding/base64"
	"fmt"
)

var (
	// symmetricKey is the key used for encrypting the configMap.
	// The same should be accessible to kata agent in order to decrypt the configMap.
	symmetricKey []byte

	// symmKeyNonce is provided by the end user. This needs to be accessible to the kata agent
	// in order to decrypt configMap.
	symmKeyNonce []byte
)

// EncryptConfigMap encrypts the configMap and returns the base64 encoded string
// of the encrypted config
func EncryptConfigMap(configMap []byte, keyPath string, noncePath string) (encConfigMap string, err error) {
	symmKey, nonce, err := GetConfigMapKeys(keyPath, noncePath)
	if err != nil {
		return "", nil
	}

	//The keys are base64 encoded. Decode it before passing to encryption function
	symmKeyDecoded, err := b64.StdEncoding.DecodeString(string(symmKey))
	if err != nil {
		fmt.Printf("Unable to decode Key\n")
		return "", err
	}
	nonceDecoded, err := b64.StdEncoding.DecodeString(string(nonce))
	if err != nil {
		fmt.Printf("Unable to decode Nonce\n")
		return "", err
	}

	block, err := aes.NewCipher(symmKeyDecoded)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(aesgcm.Seal(nil, nonceDecoded, configMap, nil)), nil
}

// GetConfigMapKeys returns the keys used for encrypting the configMap
func GetConfigMapKeys(keyPath string, noncePath string) (symmKey []byte, nonce []byte, err error) {
	if symmetricKey != nil && symmKeyNonce != nil {
		return symmetricKey, symmKeyNonce, nil
	}
	symmKey, err = getSymmetricKey(keyPath)
	if err != nil {
		return nil, nil, err
	}
	nonce, err = getNonce(noncePath)
	if err != nil {
		return nil, nil, err
	}

	return symmKey, nonce, nil
}

func getSymmetricKey(keyPath string) ([]byte, error) {
	var err error

	symmetricKey, err = readKeyFromFile(keyPath)

	return symmetricKey, err
}

func getNonce(noncePath string) ([]byte, error) {
	var err error

	symmKeyNonce, err = readKeyFromFile(noncePath)

	return symmKeyNonce, err
}

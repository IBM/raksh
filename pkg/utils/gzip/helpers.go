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

package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Extract loops through concatenated gzip file
func Extract(data []byte) [][]byte {
	var output [][]byte
	dup := make([]byte, len(data))
	copy(dup, data)
	for {
		uData, trail, err := _extractAndTrail(dup)
		if err != nil {
			return output
		}

		output = append(output, uData)
		if trail != nil {
			copy(dup, trail)
		} else {
			break
		}
	}
	return output
}

// _extractAndTrail extract the data and trail the remaining data
func _extractAndTrail(data []byte) ([]byte, []byte, error) {
	b := bytes.NewBuffer(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, nil, err
	}

	r.Multistream(false)
	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	resData := resB.Bytes()

	trailData, err := ioutil.ReadAll(b)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return resData, trailData, nil
}

// Create function creates zgip content and return in []byte
func Create(data []byte) ([]byte, error) {
	var zbuf bytes.Buffer

	zw := gzip.NewWriter(&zbuf)

	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}
	zw.Flush()

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return zbuf.Bytes(), nil
}

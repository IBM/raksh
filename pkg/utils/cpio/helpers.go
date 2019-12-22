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

package cpio

import (
	"bytes"
	"io"

	"github.com/cavaliercoder/go-cpio"
)

type File struct {
	Name string
	Body []byte
}

// Extract function will extract the bytes and return filesets
func Extract(data []byte) ([]File, error) {
	files := []File{}
	br := bytes.NewReader(data)
	r := cpio.NewReader(br)

	for {
		fbuf := new(bytes.Buffer)
		hdr, err := r.Next()
		if err == io.EOF {
			// end of cpio archive
			break
		}
		if err != nil {
			return files, err
		}
		if _, err := io.Copy(fbuf, r); err != nil {
			return files, err
		}
		f := File{
			Name: hdr.Name,
			Body: fbuf.Bytes(),
		}
		files = append(files, f)
	}
	return files, nil
}

// Create function creates cpio content and return in []byte
func Create(files []File) ([]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new cpio archive.
	w := cpio.NewWriter(buf)

	for _, file := range files {
		hdr := &cpio.Header{
			Name: file.Name,
			Mode: 0644,
			Size: int64(len(file.Body)),
		}
		if err := w.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}
	w.Flush()

	// Make sure to check the error on Close.
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

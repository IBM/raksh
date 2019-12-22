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
	"github.com/spf13/pflag"
)

var (
	Appflags *pflag.FlagSet
)

const (
	defaultScratchImage = "projectraksh/sc-scratch:latest"
)

func init() {
	Appflags = pflag.NewFlagSet("app", pflag.ExitOnError)
	Appflags.StringVarP(&filename, "filename", "f", "", "Input Kubernetes resource filename or directory")
	Appflags.StringVarP(&output, "output", "o", "", "Ouput directory, if not specified will write to the same directory with files postfixed by '-sc'")
	Appflags.StringVarP(&secureContainerImage, "image", "i", "", "SecureContainerImage resource name")
	Appflags.StringVar(&scratchImage, "scratch", defaultScratchImage, "Default scratch to be replaced in the securecontainer yamls")
}

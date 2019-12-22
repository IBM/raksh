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
	"github.com/spf13/pflag"
)

var (
	Imageflags      *pflag.FlagSet
	buildCmd        string
	baseImage       string
	initrd          string
	vmlinux         string
	push            bool
	skipAppCreation bool
)

func init() {
	Imageflags = pflag.NewFlagSet("img", pflag.ExitOnError)
	Imageflags.StringVarP(&buildCmd, "buildCmd", "c", "docker", "Specify container runtime to build the image. e.g. podman or docker")
	Imageflags.StringVarP(&baseImage, "baseImage", "b", "busybox:latest", "Name of the base image in the secure container image")
	Imageflags.StringVar(&initrd, "initrd", "", "Kata containers initrd")
	Imageflags.StringVar(&vmlinux, "vmlinux", "", "Kata containers vmlinux")
	Imageflags.BoolVarP(&push, "push", "p", false, "Push the image after building it")
	Imageflags.BoolVar(&skipAppCreation, "skip-app", false, "Skip app creation")
}

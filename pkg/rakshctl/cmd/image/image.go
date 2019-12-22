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
	"github.com/spf13/cobra"
)

var (
	imageLong = `
		Create container image with VM kernel and initrd for use with secure container `

	imageExample = `
	    rakshctl image create  -i nginx-securecontainerimage \
		--initrd /usr/share/kata-containers/kata-containers-initrd.img   \
		--vmlinux /usr/share/kata-containers/vmlinux.container \
		--symmKeyFile /home/pradipta/symmKeyFile \
		--filename /home/pradipta/nginx.yaml \
		--scratch docker-registry/sc-scratch:latest \
		--push \
		docker-registry/nginx-securecontainerimage:latest `
)

func NewCmdImage() *cobra.Command {
	cmds := &cobra.Command{
		Use:                   "image",
		DisableFlagsInUseLine: true,
		Short:                 "SecureContainer VM Image",
		Long:                  imageLong,
		Example:               imageExample,
	}
	cmds.AddCommand(NewCmdImageCreate())
	return cmds
}

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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ibm/raksh/pkg/crypto"
	appcmd "github.com/ibm/raksh/pkg/rakshctl/cmd/app"
	"github.com/ibm/raksh/pkg/rakshctl/types/flags"
	"github.com/ibm/raksh/pkg/utils/cmd"
	cpioutil "github.com/ibm/raksh/pkg/utils/cpio"
	gziputil "github.com/ibm/raksh/pkg/utils/gzip"
)

var (
	createLong = `
		Really long information goes here`

	createExample = `
		# Examples goes here`
	MissingImageError = `Image Name is missing or more than one parameters are passed`
	imageName         = "smvimage:latest"
	appArgs           []string
)

func updateImageName(image string) {
	imageName = image
}

func visitAppFlags(f *flag.Flag) {
	if f.Changed {
		appArgs = append(appArgs, "--"+f.Name, f.Value.String())
	}
}

func NewCmdImageCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [imageName]",
		DisableFlagsInUseLine: true,
		Short:                 "Create SecureContainer Image",
		Long:                  createLong,
		Example:               createExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || len(args) > 1 {
				return errors.New("Image Name is missing or more than one parameters are passed")
			}
			updateImageName(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n\nPhase 1: Executing \"kubectl image create\" command")
			fmt.Println("")
			if err := createImage(); err != nil {
				return err
			}
			if !skipAppCreation {
				appc := appcmd.NewCmdAppCreate()
				fmt.Println("\n\nPhase 2: Executing \"kubectl app create\" command")
				fmt.Println("")
				appcmd.Appflags.VisitAll(visitAppFlags)
				fmt.Println("with args: ", appArgs)
				appcmd.Appflags.Parse(appArgs)
				appc.SetArgs(appArgs)
				if err := appc.Execute(); err != nil {
					return nil
				}
			}
			return nil
		},
	}
	cmd.Flags().AddFlagSet(Imageflags)
	cmd.Flags().AddFlagSet(appcmd.Appflags)

	cmd.MarkFlagRequired("vmlinux")
	cmd.MarkFlagRequired("initrd")
	return cmd
}

func copy(source, destination string) error {

	input, err := ioutil.ReadFile(source)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		log.Println("Error creating", destination)
		log.Println(err)
		return err
	}
	return nil
}

func appendKeys(src string, dest string) error {
	initrdBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	symmKey, nonce, err := crypto.GetConfigMapKeys(flags.Key)

	var keys = []cpioutil.File{
		{Name: "symm_key", Body: symmKey},
		{Name: "nonce", Body: nonce},
	}

	cbytes, err := cpioutil.Create(keys)
	if err != nil {
		return fmt.Errorf("Failed to create cpio from keys: %+v", err)
	}

	gzbytes, err := gziputil.Create(cbytes)
	if err != nil {
		return fmt.Errorf("Failed to create gzip from cpio bytes: %+v", err)
	}

	final := append(initrdBytes, gzbytes...)
	err = ioutil.WriteFile(dest, final, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createImage() error {
	dockerfile := `FROM {{.BaseImage}}
ADD ./initrd.img /securecontainer/
ADD ./vmlinux /securecontainer/`
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir) // clean up

	dockerfn := filepath.Join(dir, "Dockerfile")
	log.Println("dockerfile: ", dockerfn)

	t := template.Must(template.New("dockerfile").Parse(dockerfile))
	data := struct {
		BaseImage string
	}{
		baseImage,
	}
	f, err := os.Create(dockerfn)
	if err != nil {
		return fmt.Errorf("Failed to a create file: %+v", err)
	}
	defer f.Close()
	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("executing template: %+v", err)
	}

	err = appendKeys(initrd, dir+"/initrd.img")
	if err != nil {
		return fmt.Errorf("Failed to append keys to given initrd: %+v", err)
	}

	err = copy(vmlinux, dir+"/vmlinux")
	if err != nil {
		return fmt.Errorf("Failed to copy %s to %s/vmlinux", vmlinux, dir)
	}

	args := []string{"build", "--no-cache", "-t", imageName, "-f", dockerfn, dir}
	_, _, err = cmd.Exec(buildCmd, args)
	if err != nil {
		return fmt.Errorf("Failed to build the SecureContainer image: %+v", err)
	}

	if push {
		_, _, err = cmd.Exec(buildCmd, []string{"push", imageName})
		if err != nil {
			return fmt.Errorf("Failed to push the %s: %+v", imageName, err)
		}
	}
	return nil
}

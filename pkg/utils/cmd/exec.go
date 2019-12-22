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

package cmd

import (
	"bufio"
	"io/ioutil"
	"log"
	"os/exec"
)

// Exec function will execute the shell command
func Exec(command string, args []string) (string, string, error) {
	var outString, errString string
	log.Printf("command: %s and args: %+v\n", command, args)
	cmd := exec.Command(command, args...)
	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error while creating the stdoutpipe for the command", err)
		return outString, errString, err
	}

	stderrpipe, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Error while creating the stderrpipe for the command", err)
		return outString, errString, err
	}

	done := make(chan struct{})

	scanner := bufio.NewScanner(stdoutpipe)
	go func() {
		for scanner.Scan() {
			st := scanner.Text()
			outString = outString + st
			log.Println(st)
		}
		done <- struct{}{}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println("Failed to start the command:", err)
		return outString, errString, err
	}

	<-done

	stderr, _ := ioutil.ReadAll(stderrpipe)
	log.Printf("%s", stderr)
	errString = string(stderr)

	err = cmd.Wait()
	if err != nil {
		log.Println("Failed to wait for the command", err)
		return outString, errString, err
	}
	return outString, errString, nil
}

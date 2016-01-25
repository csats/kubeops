// Copyright © 2016 C-SATS support@csats.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Kubeops inherited a ton of stuff from the previous set of bash scripts that kept C-SATS running. I'm in the process
// of rewriting it in Go, but in the meantime this makes it easy for me to pull those bash scripts into this.

package scripts

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var tmp string
var err error

func Run(name string, wd string, args ...string) {
	defer func() {
		os.Remove(tmp)
	}()
	var file *os.File
	file, err = ioutil.TempFile("", fmt.Sprintf("kubeops-%s-", name))
	if err != nil {
		log.Fatalf("Error creating temporary directory: %v", err)
	}
	file.Chmod(0700)
	tmp = file.Name()
	file.Close()

	data, err := Asset(name)
	if err != nil {
		log.Fatalf("Error retrieving script: %v", err)
	}
	err = ioutil.WriteFile(tmp, data, 0755)
	if err != nil {
		log.Fatalf("Error writing script: %v", err)
	}

	cmd := exec.Command(tmp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = wd

	// Fire it up!
	err = cmd.Start()
	if err != nil {
		log.Panicf("Failed to execute '%s': %s", tmp, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Script execution failed.")
	}
}

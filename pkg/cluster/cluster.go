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

package cluster

import (
	"fmt"
	"log"

	corecluster "github.com/csats/coreos-kubernetes/multi-node/aws/pkg/cluster"
)

// type Cluster struct {
// 	config corecluster.Config
// }

var FromConfig = func(file string) string {
	cfg := corecluster.NewDefaultConfig("nope")
	if err := corecluster.DecodeConfigFromFile(cfg, file); err != nil {
		log.Fatalf("Failed to decode cluster file: %s", err)
	}
	fmt.Printf("%v\n", cfg)
	return ""
}

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

package cmd

import (
	"fmt"

	"github.com/csats/kubeops/pkg/cluster"

	"github.com/csats/coreos-kubernetes/multi-node/aws/pkg/cluster"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// clusterUpCmd represents the clusterUp command
var clusterUpCmd = &cobra.Command{
	Use:   "cluster-up",
	Short: "Create a Kubernetes cluster",
	Long:  `Given a kubeops YAML file, create a CoreOS AWS Kubernetes cluster.`,
	Run:   runClusterUp,
}

func runClusterUp(cmd *cobra.Command, args []string) {
	// TODO: Work your own magic here
	fmt.Println("clusterUp called")
	fmt.Printf("%v", cluster.New)
}

func init() {
	RootCmd.AddCommand(clusterUpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterUpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterUpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

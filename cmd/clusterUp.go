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
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/csats/kubeops/pkg/cluster"
	"github.com/csats/kubeops/pkg/scripts"
)

// clusterUpCmd represents the clusterUp command
var clusterUpCmd = &cobra.Command{
	Use:   "cluster-up",
	Short: "Create a Kubernetes cluster",
	Long:  `Given a kubeops YAML file, create a CoreOS AWS Kubernetes cluster.`,
	Run:   runClusterUp,
}

// Returns true if the file exists.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func confirm() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Type 'yes' to proceed: ")
	response, _ := reader.ReadString('\n')
	if response != "yes\n" {
		os.Exit(0)
	}
}

func runClusterUp(cmd *cobra.Command, args []string) {

	c, err := cluster.FromConfig(args[0])
	if err != nil {
		fmt.Errorf("Error parsing config file: %v\n")
		os.Exit(1)
	}
	clusterName := c.Config.AWSCoreOS.ClusterName
	clusterExternalDNSName := c.Config.AWSCoreOS.ExternalDNSName
	clusterControllerIP := c.Config.AWSCoreOS.ControllerIP

	// Take care of secrets
	clusterCertDir := path.Join(secretDir, clusterName)
	exists, err := fileExists(clusterCertDir)
	if err != nil {
		fmt.Errorf("Error finding : %v\n")
		os.Exit(1)
	}
	fmt.Printf("Using secret directory %s\n", clusterCertDir)
	if exists {
		fmt.Println(`That directory appears to already exist. That's fine, it probably just means you
destroyed this cluster and now you're bringing up a new one. Proceed to overwrite
keys?`)
		confirm()
		os.RemoveAll(clusterCertDir)
	}
	os.Mkdir(clusterCertDir, 0700)
	scripts.Run("generate-keys.sh", clusterCertDir, clusterName, clusterExternalDNSName, clusterControllerIP)

	// Make the dang thing
	fmt.Printf("Creating cluster %s. Are you sure? Press enter to continue.\n", c.Config.AWSCoreOS.ClusterName)
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	if err := c.Create(); err != nil {
		fmt.Errorf("Error creating cluster: %v", err)
		os.Exit(1)
	}
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

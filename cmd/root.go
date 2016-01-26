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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csats/kubeops/pkg/cluster"
)

var cfgFile string
var secretDir string
var clusterDir string
var roleDir string
var localDir string
var localClusters []*cluster.Cluster
var kubeconfigFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubeops",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubeops.yaml)")
	RootCmd.PersistentFlags().StringVar(&secretDir, "secret-dir", "", "directory that will contain your secret keys")
	RootCmd.PersistentFlags().StringVar(&clusterDir, "cluster-dir", "./clusters", "path of your clusters directory")
	RootCmd.PersistentFlags().StringVar(&roleDir, "role-dir", "./roles", "path of your roles directory")
	RootCmd.PersistentFlags().StringVar(&localDir, "local-dir", "./.local", "path of your .local directory")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mustAbs(str string) string {
	absPath, err := filepath.Abs(str)
	if err != nil {
		fmt.Printf("Error resolving absolute path: %s\n", err)
		os.Exit(1)
	}
	return absPath
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	localDir = mustAbs(localDir)
	roleDir = mustAbs(roleDir)
	clusterDir = mustAbs(clusterDir)
	kubeconfigFile = path.Join(localDir, "kubeconfig")
	files, err := ioutil.ReadDir(clusterDir)
	if err != nil {
		fmt.Println("Error reading cluster directory. Are you in the folder of your kubeops repo?")
		fmt.Println("You can also manually specify the location with --cluster-dir=\"some-directory\"")
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	for _, f := range files {
		c, err := cluster.FromConfig(path.Join(clusterDir, f.Name()))
		if err != nil {
			// fmt.Fprintf(os.Stderr, "Error parsing %s, skipping: %s\n", f.Name(), err)
		} else {
			localClusters = append(localClusters, c)
		}
	}
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".kubeops") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

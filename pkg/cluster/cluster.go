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
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"io/ioutil"
	"path"

	corecluster "github.com/csats/coreos-kubernetes/multi-node/aws/pkg/cluster"
	"gopkg.in/yaml.v2"
)

type Cluster struct {
	Config    *Config
	TLSFiles  *TLSFiles
	SecretDir string
}

type Config struct {
	AWSCoreOS      corecluster.Config `yaml:"awsCoreOS"`
	ArtifactBucket string             `yaml:"artifactBucket"`
}

type TLSFiles struct {
	CACertFile        string
	APIServerCertFile string
	APIServerKeyFile  string
	WorkerCertFile    string
	WorkerKeyFile     string
	AdminCertFile     string
	AdminKeyFile      string
}

var FromConfig = func(file string) (*Cluster, error) {
	fmt.Printf("Decoding from %s\n", file)
	c := &Cluster{
		SecretDir: "./clusters",
		Config:    &Config{},
	}
	if err := DecodeConfigFromFile(c.Config, file); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal config file: %v", err)
	}
	c.TLSFiles = &TLSFiles{
		CACertFile:        c.getSecretPath("ca.pem"),
		APIServerCertFile: c.getSecretPath("apiserver.pem"),
		APIServerKeyFile:  c.getSecretPath("apiserver-key.pem"),
		WorkerCertFile:    c.getSecretPath("worker.pem"),
		WorkerKeyFile:     c.getSecretPath("worker-key.pem"),
		AdminCertFile:     c.getSecretPath("admin.pem"),
		AdminKeyFile:      c.getSecretPath("admin-key.pem"),
	}
	fmt.Printf("%v\n", c)
	return c, nil
}

func (cfg *Config) Valid() error {
	if cfg.ArtifactBucket == "" {
		return errors.New("artifactBucket must be set")
	}
	return cfg.AWSCoreOS.Valid()
}

func DecodeConfigFromFile(out *Config, loc string) error {
	d, err := ioutil.ReadFile(loc)
	if err != nil {
		return fmt.Errorf("failed reading config file: %v", err)
	}

	return decodeConfigBytes(out, d)
}

func decodeConfigBytes(out *Config, d []byte) error {
	if err := yaml.Unmarshal(d, &out); err != nil {
		return fmt.Errorf("failed decoding config file: %v", err)
	}

	if err := out.Valid(); err != nil {
		return fmt.Errorf("config file invalid: %v", err)
	}

	return nil
}

func (c *Cluster) getSecretPath(file string) string {
	return path.Join(c.SecretDir, file)
}

func (c *Cluster) GetStackTemplate() (string, error) {
	// Render the template
	return corecluster.StackTemplateBody(c.Config.AWSCoreOS.ArtifactURL)
}

func (c *Cluster) Create() error {
	if c.SecretDir == "" {
		return fmt.Errorf("Cluster requires SecretDir to be specified")
	}
	caCertPath := c.getSecretPath("ca.pem")
	apiserverCertPath := c.getSecretPath("apiserver.pem")
	apiserverKeyPath := c.getSecretPath("apiserver-key.pem")
	workerCertPath := c.getSecretPath("worker.pem")
	workerKeyPath := c.getSecretPath("worker-key.pem")
	adminCertPath := c.getSecretPath("admin.pem")
	adminKeyPath := c.getSecretPath("admin-key.pem")
	tlsConfig := &corecluster.TLSConfig{
		CACertFile:        caCertPath,
		CACert:            mustReadFile(caCertPath),
		APIServerCertFile: apiserverCertPath,
		APIServerCert:     mustReadFile(apiserverCertPath),
		APIServerKeyFile:  apiserverKeyPath,
		APIServerKey:      mustReadFile(apiserverKeyPath),
		WorkerCertFile:    workerCertPath,
		WorkerCert:        mustReadFile(workerCertPath),
		WorkerKeyFile:     workerKeyPath,
		WorkerKey:         mustReadFile(workerKeyPath),
		AdminCertFile:     adminCertPath,
		AdminCert:         mustReadFile(adminCertPath),
		AdminKeyFile:      adminKeyPath,
		AdminKey:          mustReadFile(adminKeyPath),
	}
	core := corecluster.New(&c.Config.AWSCoreOS, newAWSConfig(&c.Config.AWSCoreOS))
	return core.Create(tlsConfig)
}

func mustReadFile(loc string) []byte {
	b, _ := ioutil.ReadFile(loc)
	return b
}

func newAWSConfig(cfg *corecluster.Config) *aws.Config {
	c := aws.NewConfig()
	c = c.WithRegion(cfg.Region)
	return c
}

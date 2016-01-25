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
	"bytes"
	"fmt"
	"text/template"
)

var kubeconfigTemplateString = `apiVersion: v1
kind: Config
clusters:
{{range $cluster := .}}
- name: kube-aws-{{ $cluster.Config.AWSCoreOS.ClusterName }}-cluster
  cluster:
    certificate-authority: {{ $cluster.TLSFiles.CACertFile }}
    server: {{ $cluster.Config.AWSCoreOS.ExternalDNSName }}
{{end}}

contexts:
{{range $cluster := .}}
- name: kube-aws-{{ $cluster.Config.AWSCoreOS.ClusterName }}-context
  context:
    cluster: kube-aws-{{ $cluster.Config.AWSCoreOS.ClusterName }}-cluster
    namespace: default
    user: kube-aws-{{ $cluster.Config.AWSCoreOS.ClusterName }}-admin
{{end}}

users:
{{range $cluster := .}}
- name: kube-aws-{{ $cluster.Config.AWSCoreOS.ClusterName }}-admin
  user:
    client-certificate: {{ $cluster.TLSFiles.AdminCertFile }}
    client-key: {{ $cluster.TLSFiles.AdminKeyFile }}
{{end}}
current-context: broken
`

// That last line used to be
// current-context: kube-aws-{{ $cluster.ClusterName }}-context

var kubeconfigTemplate *template.Template

func init() {
	kubeconfigTemplate = template.Must(template.New("kubeconfig").Parse(kubeconfigTemplateString))
}

func GenerateKubeconfig(clusters []*Cluster) (string, error) {
	var rendered bytes.Buffer
	if err := kubeconfigTemplate.Execute(&rendered, clusters); err != nil {
		return "", fmt.Errorf("Error rendering template: %v", err)
	}
	return rendered.String(), nil
}

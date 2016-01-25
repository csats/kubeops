#!/bin/bash
# Copyright © 2016 C-SATS support@csats.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

clusterName="$1"
clusterMasterDNS="$2"
clusterControllerIP="$3"
DAYS_VALID="3650" # Change me if you don't want all the stuff to work for ten years.
(
  confFile=$(mktemp)

  #######################################################
  mildlyImportant "Generating Cluster Root CA keypair"
  #######################################################
  openssl genrsa -out ca-key.pem 2048
  openssl req \
    -x509 \
    -new \
    -nodes \
    -key ca-key.pem \
    -days 10000 \
    -out ca.pem \
    -subj "/CN=kube-ca"

  #######################################################
  mildlyImportant "Generating API Server Keypair"
  #######################################################
  cat > "$confFile" << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster.local
DNS.5 = $clusterMasterDNS
IP.1 = $clusterControllerIP
IP.2 = 10.3.0.1
EOF

  openssl genrsa -out apiserver-key.pem 2048
  openssl req \
    -new \
    -key apiserver-key.pem \
    -out apiserver.csr \
    -subj "/CN=kube-apiserver" \
    -config "$confFile"
  openssl x509 \
    -req \
    -in apiserver.csr \
    -CA ca.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out apiserver.pem \
    -days "$DAYS_VALID" \
    -extensions v3_req \
    -extfile "$confFile"

  #######################################################
  mildlyImportant "Generating Worker Keypair"
  #######################################################
  cat > "$confFile" << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = *.*.compute.internal
DNS.2 = *.ec2.internal
EOF
  openssl genrsa -out worker-key.pem 2048
  openssl req \
    -new \
    -key worker-key.pem \
    -out worker.csr \
    -subj "/CN=kube-worker" \
    -config "$confFile"
  openssl x509 \
    -req \
    -in worker.csr \
    -CA ca.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out worker.pem \
    -days "$DAYS_VALID" \
    -extensions v3_req \
    -extfile "$confFile"

  #######################################################
  mildlyImportant "Generating Admin Keypair"
  #######################################################
  cat > "$confFile" << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF
  openssl genrsa -out admin-key.pem 2048
  openssl req \
    -new \
    -key admin-key.pem \
    -out admin.csr \
    -subj "/CN=kube-admin" \
    -config "$confFile"
  openssl x509 \
    -req \
    -in admin.csr \
    -CA ca.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out admin.pem \
    -days "$DAYS_VALID" \
    -extensions v3_req \
    -extfile "$confFile"
)

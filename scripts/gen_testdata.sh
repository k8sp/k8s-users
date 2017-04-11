#!/bin/bash
set -x

root=$1
users_path=$root/users
policy_file=$root/abac-policy.jsonl

# make users directory
mkdir -p $root/users

# generate test policy files
cat > $policy_file <<EOF
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user":"*",         "nonResourcePath": "*", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user":"admin",     "namespace": "*",              "resource": "*",         "apiGroup": "*"                   }}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user":"kube-admin","namespace": "*",              "resource": "*",         "apiGroup": "*"                   }}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user":"scheduler", "namespace": "*",              "resource": "pods",                       "readonly": false }}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user":"scheduler", "namespace": "*",              "resource": "bindings"                                     }}
EOF

# generate root cert files
openssl genrsa -out $root/ca.key 2048
openssl req -x509 -new -nodes -key $root/ca.key -days 10000 -out $root/ca.crt -subj "/CN=k8s-users"

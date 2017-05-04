#!/usr/bin/env bash

# config path
CONFIG=/etc/kubernetes
# policy jsonl path
POLICY=${CONFIG}/abac/policy.jsonl
# CA cert and key path
CACRT=${CONFIG}/ssl/ca.pem
CAKEY=${CONFIG}/ssl/ca-key.pem
# SMTP host addr:port
SMTP_ADDR="smtp.partner.outlook.cn:587"
# admin email addr and secrt
EMAIL="admin@email.com"
SECRT="admin"

docker run -d -p 8091:80 \
    -v "/etc/ssl/certs:/etc/ssl/certs"\ # user for email TLS
    -v '/var/run/docker.sock:/var/run/docker.sock'\ # used for docker client to communicate to docker daemon
    -v "${CONFIG}:${CONFIG}"\
    -e ABAC_POLICY_FILE="${POLICY}"\
    -e ROOT_CA_CERT="${CACRT}"\
    -e ROOT_CA_KEY="${CAKEY}"\
    -e SMTP_SVC_ADDR=${SMTP_ADDR}\
    -e ADMIN_EMAIL=${EMAIL}\
    -e ADMIN_SECRT=${SECRT}\
    -e ADDR=":80"\
    k8s-users 

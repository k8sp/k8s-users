#!/usr/bin/env bash

# goemail package using tls to connect to smtp,
# and /etc/ssl/certs/ca-certificates.crt must exit in docker images.
# 
# -v "/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt" 

# config path
CONFIG=/etc/kubernetes
# policy jsonl path
POLICY=/abac/policy.jsonl
# CA cert and key path
CACRT=/ssl/ca.pem
CAKEY=/ssl/ca-key.pem
# SMTP host addr:port
SMTP_ADDR="smtp.partner.outlook.cn:587"
# admin email addr and secrt
EMAIL='admin@email.com'
SECRT='admin@$#\12'

docker run -d -p 8091:8091 \
    -v '/var/run/docker.sock:/var/run/docker.sock'\
    -v "${CONFIG}:${CONFIG}"\
    -e ABAC_POLICY_FILE="${POLICY}"\
    -e ROOT_CA_CERT="${CACRT}"\
    -e ROOT_CA_KEY="${CAKEY}"\
    -e SMTP_SVC_ADDR=${SMTP_ADDR}\
    -e ADMIN_EMAIL=${EMAIL}\
    -e ADMIN_SECRT=${SECRT}\
    -e ADDR=":8091" \
    bootstrapper:5000/zhanghui/k8s-users:latest



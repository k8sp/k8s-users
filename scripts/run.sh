#!/usr/bin/env bash

./k8s-users \
        -abac-policy ./testdata/abac-policy.jsonl\
        -ca-crt ./testdata/ca.crt\
        -ca-key ./testdata/ca.key \
        -smtp-svc-addr "smtp.partner.outlook.cn:587"\
        -admin-email "..." \
        -admin-secrt ""

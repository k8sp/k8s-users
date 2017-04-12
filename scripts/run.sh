#!/usr/bin/env bash

./k8s-users --ca-key ./testdata/ca.key --ca-crt ./testdata/ca.crt \
    -smtp-svc-addr 'smtp.partner.outlook.cn:587' -admin-email '...' -admin-secrt '...'

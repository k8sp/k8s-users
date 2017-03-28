#!/bin/bash
openssl genrsa -out key.pem 2048
openssl req -x509 -new -nodes -key ./key.pem -days 1000 -out ./crt.pem -subj "/CN=test"

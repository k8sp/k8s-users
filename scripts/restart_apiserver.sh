#!/bin/bash
docker restart $(docker ps |grep apiserver |awk '{print $1}')

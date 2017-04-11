#!/bin/bash
containerid=$(docker ps |grep apiserver |awk '{print $1}')
docker restart $containerid

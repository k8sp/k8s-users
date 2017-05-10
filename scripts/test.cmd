//send req to k8s-users
 curl -X POST http://localhost:8091/users -d '[
{"username":"admin","namespace":"admin","email":"zhanghui@unisound.com"},
{"username":"test","namespace":"test","email":"zhanghui@unisound.com"}
]
'

//test to host docker ps from container
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock \
                    -v /usr/bin/docker:/usr/bin/docker\
                    -v /usr/lib:/usr/lib \
                    centos:7.3.1611 /bin/bash

//target docker to restart
docker run -d ubuntu:14.04 /bin/bash -c "while ((1)); do sleep 10; done "

//worker
docker run -it -p 8091:80  -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt'\
                           -v '/var/run/docker.sock:/var/run/docker.sock'  \
                           k8s-users 


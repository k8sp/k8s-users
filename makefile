.PHONY:  build 

# generate test data 
test:
	./scripts/gen_testdata.sh ./testdata && go test ./...
# clean all
clean:
	go clean && rm -rf ./testdata/*
# local run
run:
	go build && ./scripts/run.sh
# build k8s-users docker
build:
	go build && cp k8s-users build && docker build -t bootstrapper:5000/zhanghui/k8s-users build 
# local docker run
docker-run:
	./scripts/docker-run.sh


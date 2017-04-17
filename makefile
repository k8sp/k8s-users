test:
	./scripts/gen_testdata.sh ./testdata && go test ./...
clean:
	go clean && rm -rf ./testdata/*
run:
	go build && ./scripts/run.sh
build:
	go build && cp k8s-users build && docker build -t k8s-users build 
docker-run:
	./scripts/docker-run.sh

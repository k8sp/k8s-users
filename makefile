test:
	./scripts/gen_testdata.sh ./testdata && go test ./...
clean:
	go clean && rm -rf ./testdata/*
run:
	go build && ./scripts/run.sh
build:
	go build && cp k8s-users docker && docker build -t k8s-users docker
docker-run:
	./scripts/docker-run.sh

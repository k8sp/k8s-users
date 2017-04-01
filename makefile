test:
	./scripts/gen_testdata.sh ./testdata && go test ./...
clean:
	rm -rf ./testdata/*

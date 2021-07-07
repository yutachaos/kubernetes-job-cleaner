REPOSITORY=yutachaos
REVISION =$(shell git rev-parse HEAD | head -c 8)

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

push_image:
	docker build -t ${REPOSITORY}/kubernetes-job-cleaner:$(REVISION) .
	docker push ${REPOSITORY}/kubernetes-job-cleaner:$(REVISION)
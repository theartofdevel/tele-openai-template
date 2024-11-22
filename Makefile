.PHONY: lint
lint:
	cd app; golangci-lint run -v --out-format tab --path-prefix app/

.PHONY: lint-ci
lint-ci:
	cd app; golangci-lint run -v --timeout=15m

.PHONY: run
run:
	@cd app; go build -o app cmd/main.go && ./app -config ../configs/config.yml

.PHONY: test
test:
	cd app; CGO_ENABLED=1 go test -v -race -count=1 ./...

.PHONY: build-docker
build-docker:
	@CGO_ENABLED=1 docker build -t tele-openai-bot:1.0.0 -f Dockerfile .

.PHONY: build-docker-linux
build-docker-linux:
	@GOOS=linux GOARCH=amd64 docker buildx build --platform linux/amd64 --load -t tele-openai-bot:1.0.0 -f Dockerfile .
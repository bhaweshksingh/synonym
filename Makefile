APP=synonym-service
APP_VERSION:=0.1
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

CONFIG_FILE="./.env"
HTTP_SERVE_COMMAND="http-serve"

setup: copy-config migrate

infra-local:
	docker-compose -f build/docker-compose.network.yml -f build/docker-compose.infra-basics.yml  up -d

compile:
	mkdir -p out/
	go build -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) cmd/*.go

build: deps compile

http-serve: build
	$(APP_EXECUTABLE) -configFile=$(configFile) $(HTTP_SERVE_COMMAND)

app:
	docker-compose -f build/docker-compose.app-synonym.yml -f build/docker-compose.network.yml up -d --build
	docker logs -f synonym-service-go

http-local-serve: build
	$(APP_EXECUTABLE) -configFile=$(CONFIG_FILE) $(HTTP_SERVE_COMMAND)

tidy:
	go mod tidy

deps:
	go mod download

check: fmt vet lint

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

clean:
	rm -rf out/

test:
	go clean -testcache
	go test ./...

test-cover-html:
	go clean -testcache
	mkdir -p out/
	go test ./... -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out

ci-test: test

lint:
	golangci-lint run cmd/... pkg/...

dependency-check:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore


# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./
BINARY_NAME := app

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy mod file
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v


## lint: lint code and format code
.PHONY: lint
lint: tidy
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out

## build: build the application
.PHONY: build
build: generate
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run:
	CGO_ENABLED=0 go build \
        -ldflags "-X main.gitHash=$(git rev-parse HEAD) -X main.buildStamp=$(date +%Y%m%d.%H%M%S)" \
        -o goapp main.go


## generate: go generate in the application
.PHONY: generate
generate:
	go generate ./...

## setup: tool, install deps and githook in application
.PHONY: setup
setup:
	chmod +x script/setup.sh script/setup-pre-commit.sh script/setup-dev-tool.sh script/install-deps.sh .githook/pre-commit
	./script/setup.sh

## generate-swagger
.PHONY: gen-swagger
gen-swagger:
	swag init --parseDependency --parseInternal --parseDepth 1


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push


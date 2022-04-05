export CGO_ENABLED=0
export GO111MODULE=on

# We use IMAGE_URI env var to be able to run tests for images that are stored in different 
# docker registries (google and aws). If this env var is not set, google registry is used by default.
#if test -z ${MERKELY_IMAGE_URI+y}
#then
#		export MERKELY_IMAGE_URI=ghcr.io/merkely-development/merkely:latest
#fi
[ -v ${MERKELY_IMAGE_URI} ] && export MERKELY_IMAGE_URI=ghcr.io/merkely-development/merkely:latest

LDFLAGS := -w -s

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/merkely-development/reporter/internal/version.version=${BINARY_VERSION}
endif

VERSION_METADATA = unreleased
# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X github.com/merkely-development/reporter/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X github.com/merkely-development/reporter/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/merkely-development/reporter/internal/version.gitTreeState=${GIT_DIRTY}
LDFLAGS += -extldflags "-static"

ldflags:
	@echo $(LDFLAGS)

fmt: ## Reformat package sources
	@go fmt ./...
.PHONY: fmt

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.39-alpine golangci-lint run --timeout=5m  -v ./...
.PHONY: lint

vet: fmt
	@go vet ./...
.PHONY: vet

deps: ## Install depdendencies. Runs `go get` internally.
	@GOFLAGS="" go mod download
	@GOFLAGS="" go mod tidy
.PHONY: deps

build: deps vet ## Build the binary
	@go build -o merkely -ldflags '$(LDFLAGS)' ./cmd/merkely/
.PHONY: build

test_unit: deps vet ## Run unit tests
	@docker-compose down || true
	@docker-compose up -d
	@docker exec merkely-server /demo/create_test_users.py
	@go test -v -cover -p=1 -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
.PHONY: test_unit

docker: deps vet lint
	@docker build -t merkely-cli .
.PHONY: docker

docs: build
	@rm docs.merkely.com/content/client_reference/merkely*
	@export DOCS=true && ./merkely docs --dir docs.merkely.com/content/client_reference
.PHONY: docs

licenses:
	@rm -rf licenses || true
	@go install github.com/google/go-licenses@latest
	@go-licenses save ./... --save_path="licenses/" || true
	$(eval DATA := $(shell go-licenses csv ./...))
	@echo $(DATA) | tr " " "\n" > licenses/licenses.csv
.PHONY: licenses

hugo: docs helm-docs
	cd docs.merkely.com && hugo server --minify
.PHONY: hugo

helm-lint: 
	@cd charts/k8s-reporter && helm lint .
.PHONY: helm-lint

helm-docs: helm-lint
	@cd charts/k8s-reporter &&  docker run --rm --volume "$(PWD):/helm-docs" jnorwood/helm-docs:latest --template-files README.md.gotmpl,_templates.gotmpl --output-file README.md
	@cd charts/k8s-reporter &&  docker run --rm --volume "$(PWD):/helm-docs" jnorwood/helm-docs:latest --template-files README.md.gotmpl,_templates.gotmpl --output-file ../../docs.merkely.com/content/helm/helm_chart.md
.PHONY: helm-docs

release:
	@git remote update
	@git status -uno | grep --silent "Your branch is up to date" || (echo "ERROR: your branch is NOT up to date with remote" && return 1)
	git tag -a $(tag) -m"$(tag)"
	git push origin $(tag)
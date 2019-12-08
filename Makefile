files:= $(shell find . -name '*.go' -print)
pkgs:= $(shell go list ./...)
prod_pkgs:= $(shell go list ./... | grep -v /mocks | grep -v /test)

ifdef VERSION
	version := $(VERSION)
else
	git_rev := $(shell git rev-parse --short HEAD)
	git_tag := $(shell git tag --points-at=$(git_rev))
	version := $(if $(git_tag),$(git_tag),dev-$(git_rev))
endif
build_time := $(shell date -u)
ldflags := -X "github.com/totahuanocotl/hello-world/cmd.version=$(version)" -X "github.com/totahuanocotl/hello-world/cmd.buildTime=$(build_time)"
repo := 241776843775.dkr.ecr.eu-west-1.amazonaws.com/axiltia
image :=  hello-world

.phony: all setup clean build install check checkformat format vet lint test docker-build docker-push

all : install check
check : checkformat vet lint test

setup:
	@echo "== setup"
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go mod download

clean:
	@echo "== clean"
	@go clean
	@go clean -testcache
	rm -rf build

build:
	@echo "== build"
	go build -ldflags '-s $(ldflags)'

install:
	@echo "== install"
	go install -v -ldflags '-s $(ldflags)'


unformatted:= $(shell goimports -l $(files))
checkformat:
	@echo "== check format"
ifneq "$(unformatted)" ""
	@echo "Files need formatting: $(unformatted)"
	@echo "run make format"
	@exit 1
endif

format:
	@echo "== format"
	@goimports -w $(files)
	@sync

vet:
	@echo "== vet"
	@go vet $(packages)

lint:
	@echo "== lint"
	@for pkg in $(prod_pkgs); do \
		golint -set_exit_status $$pkg || exit 1; \
	done;

test: install
	@echo "== test"
	go test -race $(pkgs)

docker-build:
	@echo "== docker-build"
	docker build \
	       -t local/$(image) \
	       -t $(repo)/$(image):$(version) \
	       .
docker-push: docker-build
	@echo "== docker-push"
	@echo "Running 'aws ecr get-login' && docker push ..."
	@LOGIN=$$(aws ecr get-login --no-include-email --region eu-west-1) && \
	       $$LOGIN && \
	       docker push $(repo)/$(image):$(version)
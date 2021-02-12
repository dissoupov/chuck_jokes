include .project/gomod-project.mk

GHE_HOST := git.soma.salesforce.com
SCRIPTS_PATH := ${PROJ_ROOT}/scripts
SHA := $(shell git rev-parse HEAD)
SHA_SHORT := $(shell git describe --always --abbrev=6)

export COVERAGE_EXCLUSIONS="vendor|tests|main.go"
export GO111MODULE=on
#BUILD_FLAGS=-mod=vendor

.PHONY: *

.SILENT:

default: all

all: vars clean folders tools generate change_log build test

clean:
	echo "Running clean"
	go clean
	rm -rf \
		./bin \
		./.rpm \
		${COVPATH} \
		/tmp/jokes/logs \
		/tmp/jokes/audit

tools:
	go install golang.org/x/tools/cmd/stringer
	go install golang.org/x/tools/cmd/gorename
	go install golang.org/x/tools/cmd/godoc
	go install golang.org/x/tools/cmd/guru
	go install golang.org/x/lint/golint
	go install github.com/go-phorce/cov-report/cmd/cov-report
	go install github.com/go-phorce/configen/cmd/configen

folders:
	mkdir -p /tmp/jokes/logs \
		/tmp/jokes/audit

version:
	echo "Building version"
	gofmt -r '"GIT_VERSION" -> "$(GIT_VERSION)"' internal/version/current.template > internal/version/current.go

build: 
	echo "*** Building jokes"
	go build ${BUILD_FLAGS} -o ${PROJ_ROOT}/bin/jokes ./cmd/jokes

change_log:
	echo "Recent changes:" > ./change_log.txt
	git log -n 20 --pretty=oneline --abbrev-commit >> ./change_log.txt

docker: build
	docker build --no-cache -t jokes -f Dockerfile .

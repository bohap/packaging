# Golang build specific envvars that would be
# propagated to build and itest stages
GO_BUILD_FOLDER=build
GO_BUILD_OUT_PATH=${GO_BUILD_FOLDER}/bin
BINARY_NAME=server
GO_TEST_OUT_PATH=${GO_BUILD_FOLDER}/testresults
GO_PKG_OUT_PATH=${GO_BUILD_FOLDER}/package
PROJECT_NAME="server"

# By default, use the host system (Mac on Mac, Linux on CI)
TARGET_OS ?= $(shell go env GOOS)
TARGET_ARCH ?= $(shell go env GOARCH)

.PHONY: help
help:
	@echo 'Usage:'
	@echo '  ## Develop / Test Commands'
	@echo '    make clean           Clean the directory tree of produced artifacts.'
	@echo '    make build           Build package'
	@echo '    make test            Run tests'
	@echo '    make testreport      Run tests and generate HTML testreport'
	@echo '    make run          	Run the app'
	@echo '    make refresh-dependencies	Update dependencies in go.mod to latest MINOR.PATCH version'

# Target clean removes the build folder specified in
# parameters and triggers golang cleanup
.PHONY: clean
clean:
	@echo "# Running clean task #"
	$(info | destroying ${GO_BUILD_FOLDER} folder)
	rm -rf ${GO_BUILD_FOLDER}
	go clean -testcache

# Build target executes golang compilation with environment
# folder given in GO_BUILD_OUT_PATH parameter
.PHONY: build
build:
	@echo "# Running build task for ${TARGET_OS}/${TARGET_ARCH} #"
	@mkdir -p ${GO_BUILD_OUT_PATH}
	GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} CGO_ENABLED=0 go build -o ${GO_BUILD_OUT_PATH}/${BINARY_NAME} .

# Test target executes golang itest execution with environment
# JSON testresult into GO_TEST_OUT_PATH
.PHONY: test
test:
	@echo "# Running test task #"
	go test ./...

# Testreport target executes golang itest execution with
# JSON testresult into GO_TEST_OUT_PATH
# then installs go-itest-report utility,
# picks JSON testresult from GO_TEST_OUT_PATH and
# generates HTML itest report
.PHONY: testreport
testreport:
	@echo "| installing go-test-report"
	go install github.com/vakenbolt/go-test-report@latest
	@echo "| creating ${GO_TEST_OUT_PATH} folder"
	mkdir -p ${GO_TEST_OUT_PATH}
	@echo "| running go test and generating report"
	go test ./... -json > ${GO_TEST_OUT_PATH}/report.json -coverprofile=${GO_TEST_OUT_PATH}/coverage.out
	cat ${GO_TEST_OUT_PATH}/report.json | go-test-report -t ${PROJECT_NAME} -o ${GO_TEST_OUT_PATH}/report.html
	go tool cover -html=${GO_TEST_OUT_PATH}/coverage.out -o ${GO_TEST_OUT_PATH}/cover.html
	@echo "| Report generated at ${GO_TEST_OUT_PATH}/report.html"

DB_HOST ?= localhost
DB_PORT ?= 6432
DB_USERNAME ?= server
DB_PASSWORD ?= server
DB_NAME ?= server

# Run the server app
.PHONY: run
run:
	@echo "# Running app on ${TARGET_OS}/${TARGET_ARCH} #"
	DB_HOST=$(DB_HOST) \
    	DB_PORT=$(DB_PORT) \
    	DB_USERNAME=$(DB_USERNAME) \
    	DB_PASSWORD=$(DB_PASSWORD) \
    	DB_NAME=$(DB_NAME) \
    	GOOS=${TARGET_OS} \
    	GOARCH=${TARGET_ARCH} \
    	go run main.go

# Refresh dependencies and update everything to
# the latest MINOR.PATCH version available
.PHONY: refresh-dependencies
refresh-dependencies:
	@echo "# Updating dependencies #"
	go get -t -u ./...
	go mod tidy

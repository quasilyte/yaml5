NOW=`date -u '+%b %d %Y %H:%M:%S'`
OS=`uname -m`
AFTER_COMMIT=`git rev-parse HEAD`
GOPATH_DIR=`go env GOPATH`

install:
	go install -ldflags "-X 'main.BuildTime=$(NOW)' -X 'main.BuildOSUname=$(OS)' -X 'main.BuildCommit=$(AFTER_COMMIT)'" ./cmd/yaml5

check:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.30.0
	@echo "running linters..."
	@$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "running tests..."
	@go test -count 3 -race -v ./...
	@echo "everything is OK"

.PHONY: check

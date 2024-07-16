.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: lint
## lint: run golangci-lint
# Install: https://golangci-lint.run/welcome/install/#local-installation
lint:
	@golangci-lint run ./... --out-format colored-line-number

.PHONY: run
## run: runs the Go program
run:
	@go run -buildvcs .

.PHONY: test
## test: runs the Go test suite
test:
	@go test -mod=readonly -v ./... -timeout=10m -coverprofile=coverage.txt

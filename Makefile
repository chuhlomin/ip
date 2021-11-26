.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs the Go program
run:
	@go run .

.PHONY: test
## test: runs the Go test suite
test:
	@go test -v ./... -timeout=10m -coverprofile=coverage.txt

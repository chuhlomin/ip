.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs the Go program
run:
	@go run -buildvcs .

.PHONY: test
## test: runs the Go test suite
test:
	@go test -mod=readonly -v ./... -timeout=10m -coverprofile=coverage.txt

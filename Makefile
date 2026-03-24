CSV_FILE=offers.csv
APP_NAME=app
CMD_PATH=./cmd/app
INPUT=input.txt
PKGS := $(shell go list ./internal/... | grep -v '/internal/app$$' | grep -v '/internal/model$$' | grep -v '/internal/logger$$')

.PHONY: run build clean test

run-mock:
	cat $(INPUT) | EE_LOG_LEVEL=debug go run $(CMD_PATH) $(CSV_FILE)

run:
	go run $(CMD_PATH) $(CSV_FILE)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

run-built:
	./$(APP_NAME) $(CSV_FILE)

test: test-unit

test-unit:
	go test -v $(PKGS) -short

test-integration:
	go test -v ./test/integration -tags=integration

test-all:
	go test -v $(PKGS) -tags=integration

coverage:
	go test $(PKGS) -coverprofile=coverage.out

coverage-html:
	go tool cover -html=coverage.out

clean:
	rm -f coverage.out
	rm -f $(APP_NAME)
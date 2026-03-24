CSV_FILE=offers.csv
APP_NAME=app
CMD_PATH=./cmd/app
INPUT=input.txt

.PHONY: run build clean test

run-mock:
	EE_LOG_LEVEL=debug go run $(CMD_PATH) $(CSV_FILE) < $(INPUT)

run:
	go run $(CMD_PATH) $(CSV_FILE)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

run-built:
	./$(APP_NAME) $(CSV_FILE)

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)
APP_NAME=app
CMD_PATH=./cmd/app
INPUT=input.txt

.PHONY: run build clean test

run:
	go run $(CMD_PATH) < $(INPUT)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

run-built:
	./$(APP_NAME) < $(INPUT)

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)
TARGET_DIR = build
TARGET = $(TARGET_DIR)/main

all: run

$(TARGET): cmd/main.go
	go build -o $@ $^ 

build: $(TARGET)

run: build
	./$(TARGET)

test:
	go test -v ./...

clean:
	@rm -rf $(TARGET_DIR) 

.PHONY: all build clean test run

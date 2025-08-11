.PHONY: build test run clean lint vet

BINARY_NAME=wb-calendar
BUILD_DIR=build
MAIN_FILE=cmd/main.go

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	go test -v ./...

run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_FILE)

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f calendar.log
	@echo "Clean completed"
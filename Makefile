BINARY_NAME := ytcw
BUILD_DIR := build
SRC := main.go

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

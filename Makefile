BINARY_NAME := ytcw
BUILD_DIR := build
SRC := main.go
SWAG_CMD := swag init -g ./cmd/serve.go
SWAG_DIR := docs

.PHONY: all
all: docs build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: docs
docs:
	@echo "Generating Swagger docs..."
	@$(SWAG_CMD)
	@echo "Swagger docs generated successfully"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -rf $(BUILD_DIR)
	@rm -rf $(SWAG_DIR)
	@echo "Clean complete"

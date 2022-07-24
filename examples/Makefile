BUILD=build

all: clean build

build: go.sum
	@echo "Building ..."
	@go build -o $(BUILD)/go-embedding

go.sum: go.mod
	@echo "Ensure dependencies have not been modified"
	@go mod tidy
	@GO111MODULE=on go mod verify

clean:
	@echo "Clean old built"
	@rm -rf $(BUILD)
	@go clean
